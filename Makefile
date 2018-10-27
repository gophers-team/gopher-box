.PHONY: build build_device build_main build_main_local deploy deploy_device

build_device:
	GOARCH=arm64 GOOS=linux go build -o ./build/device ./device

build_main:
	GOARCH=amd64 GOOS=linux go build -o ./build/gopher-box ./server

build: build_device build_main

build_main_local:
	go build -o ./build/gopher-box ./server

deploy: build_main
	scp ./gopher-box.service root@130.193.56.206:/etc/systemd/system/gopher-box.service
	scp ./build/gopher-box root@130.193.56.206:/usr/local/bin/gopher-box
	ssh root@130.193.56.206 'systemctl stop gopher-box.service; systemctl daemon-reload; systemctl restart gopher-box.service'

deploy_device: build_device
	scp -C ./build/device linaro@172.31.19.157:
