Gopher Box aka CureBox
-----

CureBox is a smart solution which makes scheduled medication simple and controllable. It allows patient to simplify following medication schedule and connect with doctor and relatives. Integrated admission monitoring also enables notifying loved ones in case of medication schedule violation.

[Presentation](https://www.dropbox.com/s/fcrtbxbanojdvke/CureBox_by_Gophers.pdf)

# Quickstart

You will need [golang](https://golang.org/) compiler with [modules](https://github.com/golang/go/wiki/Modules) support to build the binaries.

Build Server and run it on your host:
```
make build_main_local
./build/gopher-box -db-file ./events.db -init-db -port 8080 -server
```

Build Device for arm64:
```
make build_device
./build/device
```

Or build Device for local use:
```
make build_device_local
./build/device
```

Device requires your Arduino to be flahed with [Firmata Firmware](https://github.com/firmata/arduino). You can make it with [gort](http://gort.io/), gobot has [instructions](https://gobot.io/documentation/platforms/arduino/) for different platforms.