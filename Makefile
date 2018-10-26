

build_main:
	GOARCH=amd64 GOOS=linux go build -o gopher-box main.go


deploy:
	scp gopher-box 130.193.56.206:/tmp/gopher-box

