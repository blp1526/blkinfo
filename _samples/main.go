package main

import (
	"fmt"

	"github.com/blp1526/go-udevinfo/blk"
)

func main() {
	mountpoint := "/"
	fmt.Printf("## mountpoint\n%s\n\n", mountpoint)

	devPath, _ := blk.GetDevPath("/")
	fmt.Printf("## devPath\n%s\n\n", devPath)

	majorMinor, _ := blk.GetMajorMinor(devPath)
	fmt.Printf("## majorMinor\n%s\n\n", majorMinor)

	udevData, _ := blk.GetUdevData(majorMinor)
	fmt.Printf("## udevData\n%s\n", udevData)
}
