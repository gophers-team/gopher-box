.PHONY: build build_device build_main build_main_local deploy deploy_device

build_device:
	GOARCH=arm64 GOOS=linux go build -o ./build/device ./device

build_main:
	GOARCH=amd64 GOOS=linux go build -o ./build/gopher-box ./server

build: build_device build_main

build_main_local:
	go build -o ./build/gopher-box ./server

deploy:
	scp go.mod root@130.193.56.206:/srv/go.mod
	scp go.sum root@130.193.56.206:/srv/go.sum
	scp -r ./server root@130.193.56.206:/srv
	scp -r ./api root@130.193.56.206:/srv
	scp ./gopher-box.service root@130.193.56.206:/etc/systemd/system/gopher-box.service
	scp -r ./static root@130.193.56.206:/srv/static
	ssh root@130.193.56.206 'systemctl stop gopher-box.service;'
	ssh root@130.193.56.206 'cd /srv && /usr/local/go/bin/go build -o /usr/local/bin/gopher-box ./server'
	ssh root@130.193.56.206 'systemctl restart gopher-box.service; systemctl daemon-reload; systemctl restart gopher-box.service'

deploy_static:
	scp -r ./static root@130.193.56.206:/srv

deploy_device: build_device
	scp -C ./build/device linaro@172.31.19.157:
