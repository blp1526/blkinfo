[![Build Status](https://travis-ci.org/blp1526/blkinfo.svg?branch=master)](https://travis-ci.org/blp1526/blkinfo)
[![Go Report Card](https://goreportcard.com/badge/github.com/blp1526/blkinfo)](https://goreportcard.com/report/github.com/blp1526/blkinfo)
[![GoDoc](https://godoc.org/github.com/blp1526/blkinfo?status.svg)](https://godoc.org/github.com/blp1526/blkinfo)
[![GolangCI](https://golangci.com/badges/github.com/blp1526/blkinfo.svg)](https://golangci.com/r/github.com/blp1526/blkinfo)

# blkinfo

A Linux Block Device Info Library

## Installation

```
$ wget https://github.com/blp1526/blkinfo/releases/latest/download/blkinfo_linux_x86_64.tar.gz
$ tar zxvf blkinfo_linux_x86_64.tar.gz
```

## Usage

```
$ blkinfo /dev/vda3
```

```json
{
  "path": "/dev/vda3",
  "resolved_path": "/dev/vda3",
  "parent_path": "/dev/vda",
  "child_paths": [],
  "sys_path": "/sys/block/vda/vda3",
  "resolved_sys_path": "/sys/devices/pci0000:00/0000:00:05.0/virtio2/block/vda/vda3",
  "sys": {
    "uevent": [
      "MAJOR=252",
      "MINOR=3",
      "DEVNAME=vda3",
      "DEVTYPE=partition",
      "PARTN=3"
    ],
    "slaves": [],
    "holders": []
  },
  "major_minor": "252:3",
  "udev_data_path": "/run/udev/data/b252:3",
  "udev_data": [
    "S:disk/by-uuid/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "S:disk/by-partuuid/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "S:disk/by-path/virtio-pci-0000:00:05.0-part3",
    "S:disk/by-path/pci-0000:00:05.0-part3",
    "W:4",
    "I:1583813",
    "E:ID_SCSI=1",
    "E:ID_PART_TABLE_TYPE=gpt",
    "E:ID_PART_TABLE_UUID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "E:ID_PATH=pci-0000:00:05.0",
    "E:ID_PATH_TAG=pci-0000_00_05_0",
    "E:ID_FS_UUID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "E:ID_FS_UUID_ENC=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "E:ID_FS_VERSION=1.0",
    "E:ID_FS_TYPE=ext4",
    "E:ID_FS_USAGE=filesystem",
    "E:ID_PART_ENTRY_SCHEME=gpt",
    "E:ID_PART_ENTRY_UUID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "E:ID_PART_ENTRY_TYPE=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "E:ID_PART_ENTRY_NUMBER=3",
    "E:ID_PART_ENTRY_OFFSET=8392704",
    "E:ID_PART_ENTRY_SIZE=33548288",
    "E:ID_PART_ENTRY_DISK=252:0",
    "E:net.ifnames=0",
    "G:systemd"
  ],
  "mount_info_path": "/proc/self/mountinfo",
  "mount_info": {
    "mount_id": "28",
    "parent_id": "0",
    "major_minor": "252:3",
    "root": "/",
    "mount_point": "/",
    "mount_options": [
      "rw",
      "relatime"
    ],
    "optional_fields": [
      "shared:1"
    ],
    "filesystem_type": "ext4",
    "mount_source": "/dev/vda3",
    "super_options": [
      "rw",
      "errors=remount-ro",
      "data=ordered"
    ]
  },
  "os_release_path": "/etc/os-release",
  "os_release": {
    "BUG_REPORT_URL": "https://bugs.launchpad.net/ubuntu/",
    "HOME_URL": "https://www.ubuntu.com/",
    "ID": "ubuntu",
    "ID_LIKE": "debian",
    "NAME": "Ubuntu",
    "PRETTY_NAME": "Ubuntu 18.04.3 LTS",
    "PRIVACY_POLICY_URL": "https://www.ubuntu.com/legal/terms-and-policies/privacy-policy",
    "SUPPORT_URL": "https://help.ubuntu.com/",
    "UBUNTU_CODENAME": "bionic",
    "VERSION": "18.04.3 LTS (Bionic Beaver)",
    "VERSION_CODENAME": "bionic",
    "VERSION_ID": "18.04"
  }
}
```

### Options

|Option|Description|
|---|---|
|--help, -h|show help|
|--output value, -o value|output as "json" or "yaml" (default: "json")|
|--version, -v|print the version|

## Package

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/blp1526/blkinfo"
)

func main() {
	bi, _ := blkinfo.New("/dev/vda3")
	b, _ := json.Marshal(bi)
	fmt.Println(string(b))
}
```

## Build

```
$ make
```
