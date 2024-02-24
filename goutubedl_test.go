package goutubedl_test

// TODO: currently the tests only run on linux as they use osleaktest which only
// has linux support

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/wader/goutubedl"
	"github.com/wader/osleaktest"
)

const (
	testVideoRawURL          = "https://www.youtube.com/watch?v=C0DPdy98e4c"
	playlistRawURL           = "https://soundcloud.com/mattheis/sets/kindred-phenomena"
	channelRawURL            = "https://www.youtube.com/channel/UCHDm-DKoMyJxKVgwGmuTaQA"
	subtitlesTestVideoRawURL = "https://www.youtube.com/watch?v=QRS8MkLhQmM"
)

func leakChecks(t *testing.T) func() {
	leakFn := leaktest.Check(t)
	osLeakFn := osleaktest.Check(t)

	return func() {
		leakFn()
		osLeakFn()
	}
}

func TestBinaryNotPath(t *testing.T) {
	defer leakChecks(t)()
	defer func(orig string) { goutubedl.Path = orig }(goutubedl.Path)
	goutubedl.Path = "/non-existing"

	_, versionErr := goutubedl.Version(context.Background())
	if versionErr == nil || !strings.Contains(versionErr.Error(), "no such file or directory") {
		t.Fatalf("err should be nil 'no such file or directory': %v", versionErr)
	}
}

func TestVersion(t *testing.T) {
	defer leakChecks(t)()

	versionRe := regexp.MustCompile(`^\d{4}\.\d{2}.\d{2}.*$`)
	version, versionErr := goutubedl.Version(context.Background())

	if versionErr != nil {
		t.Fatalf("err: %s", versionErr)
	}

	if !versionRe.MatchString(version) {
		t.Errorf("version %q does not match %q", version, versionRe)
	}
}

func TestDownload(t *testing.T) {
	defer leakChecks(t)()

	stderrBuf := &bytes.Buffer{}
	r, err := goutubedl.New(context.Background(), testVideoRawURL, goutubedl.Options{
		StderrFn: func(cmd *exec.Cmd) io.Writer {
			return stderrBuf
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	dr, err := r.Download(context.Background(), r.Info.Formats[0].FormatID)
	if err != nil {
		t.Fatal(err)
	}
	downloadBuf := &bytes.Buffer{}
	n, err := io.Copy(downloadBuf, dr)
	if err != nil {
		t.Fatal(err)
	}
	dr.Close()

	if n != int64(downloadBuf.Len()) {
		t.Errorf("copy n not equal to download buffer: %d!=%d", n, downloadBuf.Len())
	}

	if n < 10000 {
		t.Errorf("should have copied at least 10000 bytes: %d", n)
	}

	if !strings.Contains(stderrBuf.String(), "Destination") {
		t.Errorf("did not find expected log message on stderr: %q", stderrBuf.String())
	}
}

func TestDownloadWithoutInfo(t *testing.T) {
	defer leakChecks(t)()

	stderrBuf := &bytes.Buffer{}
	dr, err := goutubedl.Download(context.Background(), testVideoRawURL, goutubedl.Options{
		StderrFn: func(cmd *exec.Cmd) io.Writer {
			return stderrBuf
		},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	downloadBuf := &bytes.Buffer{}
	n, err := io.Copy(downloadBuf, dr)
	if err != nil {
		t.Fatal(err)
	}
	dr.Close()

	if n != int64(downloadBuf.Len()) {
		t.Errorf("copy n not equal to download buffer: %d!=%d", n, downloadBuf.Len())
	}

	if n < 10000 {
		t.Errorf("should have copied at least 10000 bytes: %d", n)
	}

	if !strings.Contains(stderrBuf.String(), "Destination") {
		t.Errorf("did not find expected log message on stderr: %q", stderrBuf.String())
	}
}

func TestParseInfo(t *testing.T) {
	for _, c := range []struct {
		url           string
		expectedTitle string
	}{
		{"https://soundcloud.com/avalonemerson/avalon-emerson-live-at-printworks-london-march-2017", "Avalon Emerson Live at Printworks London 2017"},
		{"https://www.infoq.com/presentations/Simple-Made-Easy", "Simple Made Easy - InfoQ"},
		{"https://www.youtube.com/watch?v=uVYWQJ5BB_w", "A Radiolab Producer on the Making of a Podcast"},
	} {
		t.Run(c.url, func(t *testing.T) {
			defer leakChecks(t)()

			ctx, cancelFn := context.WithCancel(context.Background())
			ydlResult, err := goutubedl.New(ctx, c.url, goutubedl.Options{
				DownloadThumbnail: true,
			})
			if err != nil {
				cancelFn()
				t.Errorf("failed to parse: %v", err)
				return
			}
			cancelFn()

			yi := ydlResult.Info
			results := ydlResult.Formats()

			if yi.Title != c.expectedTitle {
				t.Errorf("expected title %q got %q", c.expectedTitle, yi.Title)
			}

			if yi.Thumbnail != "" && len(yi.ThumbnailBytes) == 0 {
				t.Errorf("expected thumbnail bytes")
			}

			var dummy map[string]interface{}
			if err := json.Unmarshal(ydlResult.RawJSON, &dummy); err != nil {
				t.Errorf("failed to parse RawJSON")
			}

			if len(results) == 0 {
				t.Errorf("expected formats")
			}

			for _, f := range results {
				if f.FormatID == "" {
					t.Errorf("expected to have FormatID")
				}
				if f.Ext == "" {
					t.Errorf("expected to have Ext")
				}
				if (f.ACodec == "" || f.ACodec == "none") &&
					(f.VCodec == "" || f.VCodec == "none") &&
					f.Ext == "" {
					t.Errorf("expected to have some media: audio %q video %q ext %q", f.ACodec, f.VCodec, f.Ext)
				}
			}
		})
	}
}

func TestPlaylist(t *testing.T) {
	defer leakChecks(t)()

	ydlResult, ydlResultErr := goutubedl.New(context.Background(), playlistRawURL, goutubedl.Options{
		Type:              goutubedl.TypePlaylist,
		DownloadThumbnail: false,
	})

	if ydlResultErr != nil {
		t.Errorf("failed to download: %s", ydlResultErr)
	}

	expectedTitle := "Kindred Phenomena"
	if ydlResult.Info.Title != expectedTitle {
		t.Errorf("expected title %q got %q", expectedTitle, ydlResult.Info.Title)
	}

	expectedEntries := 8
	if len(ydlResult.Info.Entries) != expectedEntries {
		t.Errorf("expected %d entries got %d", expectedEntries, len(ydlResult.Info.Entries))
	}

	expectedTitleOne := "A1 Mattheis - Herds"
	if ydlResult.Info.Entries[0].Title != expectedTitleOne {
		t.Errorf("expected title %q got %q", expectedTitleOne, ydlResult.Info.Entries[0].Title)
	}
}

func TestChannel(t *testing.T) {
	defer leakChecks(t)()

	ydlResult, ydlResultErr := goutubedl.New(
		context.Background(),
		channelRawURL,
		goutubedl.Options{
			Type:              goutubedl.TypeChannel,
			DownloadThumbnail: false,
		},
	)

	if ydlResultErr != nil {
		t.Errorf("failed to download: %s", ydlResultErr)
	}

	expectedTitle := "Simon Yapp"
	if ydlResult.Info.Title != expectedTitle {
		t.Errorf("expected title %q got %q", expectedTitle, ydlResult.Info.Title)
	}

	expectedEntries := 5
	if len(ydlResult.Info.Entries) != expectedEntries {
		t.Errorf("expected %d entries got %d", expectedEntries, len(ydlResult.Info.Entries))
	}

	expectedTitleOne := "#RNLI Shoreham #LifeBoat demo of launch."
	if ydlResult.Info.Entries[0].Title != expectedTitleOne {
		t.Errorf("expected title %q got %q", expectedTitleOne, ydlResult.Info.Entries[0].Title)
	}
}

func TestUnsupportedURL(t *testing.T) {
	defer leaktest.Check(t)()

	_, ydlResultErr := goutubedl.New(context.Background(), "https://www.google.com", goutubedl.Options{})
	if ydlResultErr == nil {
		t.Errorf("expected unsupported url")
	}
	expectedErrPrefix := "Unsupported URL:"
	if ydlResultErr != nil && !strings.HasPrefix(ydlResultErr.Error(), expectedErrPrefix) {
		t.Errorf("expected error prefix %q got %q", expectedErrPrefix, ydlResultErr.Error())

	}
}

func TestPlaylistWithPrivateVideo(t *testing.T) {
	defer leaktest.Check(t)()

	playlistRawURL := "https://www.youtube.com/playlist?list=PLX0g748fkegS54oiDN4AXKl7BR7mLIydP"
	ydlResult, ydlResultErr := goutubedl.New(context.Background(), playlistRawURL, goutubedl.Options{
		Type:              goutubedl.TypePlaylist,
		DownloadThumbnail: false,
	})

	if ydlResultErr != nil {
		t.Errorf("failed to download: %s", ydlResultErr)
	}

	expectedLen := 2
	actualLen := len(ydlResult.Info.Entries)
	if expectedLen != actualLen {
		t.Errorf("expected len %d got %d", expectedLen, actualLen)
	}
}

func TestSubtitles(t *testing.T) {
	defer leakChecks(t)()

	ydlResult, ydlResultErr := goutubedl.New(
		context.Background(),
		subtitlesTestVideoRawURL,
		goutubedl.Options{
			DownloadSubtitles: true,
		})

	if ydlResultErr != nil {
		t.Errorf("failed to download: %s", ydlResultErr)
	}

	for _, subtitles := range ydlResult.Info.Subtitles {
		for _, subtitle := range subtitles {
			if subtitle.Ext == "" {
				t.Errorf("%s: %s: expected extension", ydlResult.Info.URL, subtitle.Language)
			}
			if subtitle.Language == "" {
				t.Errorf("%s: %s: expected language", ydlResult.Info.URL, subtitle.Language)
			}
			if subtitle.URL == "" {
				t.Errorf("%s: %s: expected url", ydlResult.Info.URL, subtitle.Language)
			}
			if len(subtitle.Bytes) == 0 {
				t.Errorf("%s: %s: expected bytes", ydlResult.Info.URL, subtitle.Language)
			}
		}
	}
}

func TestDownloadSections(t *testing.T) {
	defer leakChecks(t)()

	fileName := "durationTestingFile"
	duration := 5

	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("failed to check ffmpeg installed: %s", err)
	}

	ydlResult, ydlResultErr := goutubedl.New(
		context.Background(),
		"https://www.youtube.com/watch?v=OyuL5biOQ94",
		goutubedl.Options{
			DownloadSections: fmt.Sprintf("*0:0-0:%d", duration),
		})

	if ydlResult.Options.DownloadSections != "*0:0-0:5" {
		t.Errorf("failed to setup --download-sections")
	}
	if ydlResultErr != nil {
		t.Errorf("failed to download: %s", ydlResultErr)
	}
	dr, err := ydlResult.Download(context.Background(), "best")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = io.Copy(f, dr)
	if err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("ffprobe", "-v", "quiet", "-show_entries", "format=duration", fileName)
	stdout, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	var gotDurationString string
	output := string(stdout)
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "duration") {
			if d, found := strings.CutPrefix(line, "duration="); found {
				gotDurationString = d
			}
		}
	}

	gotDuration, err := strconv.ParseFloat(gotDurationString, 32)
	if err != nil {
		t.Fatal(err)
	}
	seconds := int(gotDuration)
	if seconds != duration {
		t.Fatalf("didnot get expected duration of %d, but got %d", duration, seconds)
	}
	dr.Close()
}

func TestErrorNotAPlaylist(t *testing.T) {
	defer leakChecks(t)()

	_, ydlResultErr := goutubedl.New(context.Background(), testVideoRawURL, goutubedl.Options{
		Type:              goutubedl.TypePlaylist,
		DownloadThumbnail: false,
	})
	if ydlResultErr != goutubedl.ErrNotAPlaylist {
		t.Errorf("expected is playlist error, got %s", ydlResultErr)
	}
}

func TestErrorNotASingleEntry(t *testing.T) {
	defer leakChecks(t)()

	_, ydlResultErr := goutubedl.New(context.Background(), playlistRawURL, goutubedl.Options{
		Type:              goutubedl.TypeSingle,
		DownloadThumbnail: false,
	})
	if ydlResultErr != goutubedl.ErrNotASingleEntry {
		t.Fatalf("expected is single entry error, got %s", ydlResultErr)
	}
}

func TestOptionDownloader(t *testing.T) {
	defer leakChecks(t)()

	ydlResult, ydlResultErr := goutubedl.New(
		context.Background(),
		testVideoRawURL,
		goutubedl.Options{
			Downloader: "ffmpeg",
		})

	if ydlResultErr != nil {
		t.Fatalf("failed to download: %s", ydlResultErr)
	}
	dr, err := ydlResult.Download(context.Background(), ydlResult.Info.Formats[0].FormatID)
	if err != nil {
		t.Fatal(err)
	}
	downloadBuf := &bytes.Buffer{}
	_, err = io.Copy(downloadBuf, dr)
	if err != nil {
		t.Fatal(err)
	}
	dr.Close()
}

func TestInvalidOptionTypeField(t *testing.T) {
	defer leakChecks(t)()

	_, err := goutubedl.New(context.Background(), playlistRawURL, goutubedl.Options{
		Type: 42,
	})
	if err == nil {
		t.Error("should have failed")
	}
}

func TestDownloadPlaylistEntry(t *testing.T) {
	defer leakChecks(t)()
	// Download file by specifying the playlist index
	stderrBuf := &bytes.Buffer{}
	r, err := goutubedl.New(context.Background(), playlistRawURL, goutubedl.Options{
		StderrFn: func(cmd *exec.Cmd) io.Writer {
			return stderrBuf
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedTitle := "Kindred Phenomena"
	if r.Info.Title != expectedTitle {
		t.Errorf("expected title %q got %q", expectedTitle, r.Info.Title)
	}

	expectedEntries := 8
	if len(r.Info.Entries) != expectedEntries {
		t.Errorf("expected %d entries got %d", expectedEntries, len(r.Info.Entries))
	}

	expectedTitleOne := "B1 Mattheis - Ben M"
	playlistIndex := 2
	if r.Info.Entries[playlistIndex].Title != expectedTitleOne {
		t.Errorf("expected title %q got %q", expectedTitleOne, r.Info.Entries[playlistIndex].Title)
	}

	dr, err := r.DownloadWithOptions(context.Background(), goutubedl.DownloadOptions{
		PlaylistIndex: int(r.Info.Entries[playlistIndex].PlaylistIndex),
		Filter:        r.Info.Entries[playlistIndex].Formats[0].FormatID,
	})
	if err != nil {
		t.Fatal(err)
	}
	playlistBuf := &bytes.Buffer{}
	n, err := io.Copy(playlistBuf, dr)
	if err != nil {
		t.Fatal(err)
	}
	dr.Close()

	if n != int64(playlistBuf.Len()) {
		t.Errorf("copy n not equal to download buffer: %d!=%d", n, playlistBuf.Len())
	}

	if n < 10000 {
		t.Errorf("should have copied at least 10000 bytes: %d", n)
	}

	if !strings.Contains(stderrBuf.String(), "Destination") {
		t.Errorf("did not find expected log message on stderr: %q", stderrBuf.String())
	}

	// Download the same file but with the direct link
	url := "https://soundcloud.com/mattheis/b1-mattheis-ben-m"
	stderrBuf = &bytes.Buffer{}
	r, err = goutubedl.New(context.Background(), url, goutubedl.Options{
		StderrFn: func(cmd *exec.Cmd) io.Writer {
			return stderrBuf
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if r.Info.Title != expectedTitleOne {
		t.Errorf("expected title %q got %q", expectedTitleOne, r.Info.Title)
	}

	expectedEntries = 0
	if len(r.Info.Entries) != expectedEntries {
		t.Errorf("expected %d entries got %d", expectedEntries, len(r.Info.Entries))
	}

	dr, err = r.Download(context.Background(), r.Info.Formats[0].FormatID)
	if err != nil {
		t.Fatal(err)
	}
	directLinkBuf := &bytes.Buffer{}
	n, err = io.Copy(directLinkBuf, dr)
	if err != nil {
		t.Fatal(err)
	}
	dr.Close()

	if n != int64(directLinkBuf.Len()) {
		t.Errorf("copy n not equal to download buffer: %d!=%d", n, directLinkBuf.Len())
	}

	if n < 10000 {
		t.Errorf("should have copied at least 10000 bytes: %d", n)
	}

	if !strings.Contains(stderrBuf.String(), "Destination") {
		t.Errorf("did not find expected log message on stderr: %q", stderrBuf.String())
	}

	if directLinkBuf.Len() != playlistBuf.Len() {
		t.Errorf("not the same content size between the playlist index entry and the direct link entry: %d != %d", playlistBuf.Len(), directLinkBuf.Len())
	}

	if !bytes.Equal(directLinkBuf.Bytes(), playlistBuf.Bytes()) {
		t.Error("not the same content between the playlist index entry and the direct link entry")
	}
}
