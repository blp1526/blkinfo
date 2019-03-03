[![Build Status](https://travis-ci.org/blp1526/blkinfo.svg?branch=master)](https://travis-ci.org/blp1526/blkinfo)

# blkinfo

A Linux Block Device Info Library

## Usage

### As A CLI

```
$ make
$ ./bin/blkinfo show --format json /dev/vda
```

```json
{
  "path": "/dev/vda",
  "real_path": "/dev/vda",
  "mountpoint": "",
  "parent_path": "",
  "child_paths": [
    "/dev/vda1",
    "/dev/vda2",
    "/dev/vda5"
  ],
  "sys_path": "/sys/devices/pci0000:00/0000:00:07.0/virtio2/block/vda",
  "sys": {
    "uevent": [
      "MAJOR=253",
      "MINOR=0",
      "DEVNAME=vda",
      "DEVTYPE=disk"
    ],
    "slaves": [],
    "holders": []
  },
  "udev_data_path": "/run/udev/data/b253:0",
  "udev_data": [
    "S:disk/by-path/virtio-pci-0000:00:07.0",
    "W:12",
    "I:2277684",
    "E:ID_PATH=virtio-pci-0000:00:07.0",
    "E:ID_PATH_TAG=virtio-pci-0000_00_07_0",
    "E:ID_PART_TABLE_UUID=xxxxxxxx",
    "E:ID_PART_TABLE_TYPE=dos",
    "E:ID_FS_TYPE=",
    "G:systemd"
  ]
}
```

### As A Package

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/blp1526/blkinfo"
)

func main() {
	bi, _ := blkinfo.New("/dev/vda")
	b, _ := json.Marshal(bi)
	fmt.Println(string(b))
}
```
