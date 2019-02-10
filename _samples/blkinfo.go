package main

import (
	"fmt"
	"strings"

	"github.com/blp1526/go-blkinfo"
)

func main() {
	blkInfo, err := blkinfo.New("/dev/sda1")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", blkInfo)

	udevInfo, err := blkInfo.UdevInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n=== blkInfo.UdevInfo() ===\n")
	for _, line := range strings.Split(udevInfo, "\n") {
		fmt.Println(line)
	}

	osInfo, err := blkInfo.OsInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n=== blkInfo.OsInfo() ===\n")
	for _, line := range strings.Split(osInfo, "\n") {
		fmt.Println(line)
	}
}
