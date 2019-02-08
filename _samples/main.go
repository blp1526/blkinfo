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
	swapInfo, err := blk.SwapInfo()
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## swapinfo\n%s\n\n", swapInfo)

	mountpoint := "/"
	fmt.Printf("## mountpoint\n%s\n\n", mountpoint)

	osType, err := blk.GetOsType(mountpoint)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## osType\n%s\n\n", osType)

	devPath, err := blk.GetDevPath(mountpoint)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## devPath\n%s\n\n", devPath)

	majorMinor, err := blk.GetMajorMinor(devPath)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## majorMinor\n%s\n\n", majorMinor)

	udevData, err := blk.GetUdevData(majorMinor)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## udevData\n%s\n\n", udevData)

	partTableType, _ := blk.GetPartTableType(majorMinor)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## partTableType\n%s\n\n", partTableType)

	fmt.Printf("###############################\n\n")

	loopDevPath := "/dev/loop0"
	loopMountpoint, err := blk.GetMountpoint(loopDevPath)
	if err != nil {
		exit1(err)
	}

	fmt.Printf("## loopDevPath: %s\nloopMountpoint: %s\n\n", loopDevPath, loopMountpoint)
}
