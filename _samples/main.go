package main

import (
	"fmt"
	"os"

	"github.com/blp1526/go-udevinfo/blk"
)

func exit1(err error) {
	fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
	os.Exit(1)
}

func main() {
	mountpoint := "/"
	fmt.Printf("## mountpoint\n%s\n\n", mountpoint)

	devPath, err := blk.GetDevPath(mountpoint)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## devPath\n%s\n\n", devPath)

	majorMinor, _ := blk.GetMajorMinor(devPath)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## majorMinor\n%s\n\n", majorMinor)

	udevData, _ := blk.GetUdevData(majorMinor)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## udevData\n%s\n", udevData)
}
