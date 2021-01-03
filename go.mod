module github.com/wader/goutubedl

go 1.12

require (
	// bump: leaktest /github.com\/fortytw2\/leaktest v(.*)/ git:https://github.com/fortytw2/leaktest.git|^1
	// bump: leaktest command go get github.com/fortytw2/leaktest@v$LATEST && go mod tidy
	github.com/fortytw2/leaktest v1.3.0
	// bump: osleaktest /github.com\/wader\/osleaktest .*-(.*)/ gitrefs:https://github.com/wader/osleaktest.git|re:%refs/heads/master%|@commit|/^(.{12})/
	// bump: osleaktest command go get github.com/wader/osleaktest && go mod tidy
	github.com/wader/osleaktest v0.0.0-20191111175233-f643b0fed071
)
