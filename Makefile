.PHONY: build build_device build_main deploy

build_device:
	GOARCH=arm64 GOOS=linux go build -o ./build/device ./device


build_main:
	GOARCH=amd64 GOOS=linux go build -o ./build/gopher-box ./server

build: build_device build_main

build_main_local:
	go build -o ./build/gopher-box ./server

deploy: build_main
	scp ./scripts/deploy.sh 130.193.56.206:/tmp/deploy.sh
	scp ./build/gopher-box 130.193.56.206:/tmp/gopher-box

