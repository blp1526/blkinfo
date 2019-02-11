package main

import (
	blkinfo "github.com/blp1526/go-blkinfo"
	"github.com/k0kubun/pp"
)

func main() {
	blkInfo, err := blkinfo.New("/dev/sda1")
	if err != nil {
		panic(err)
	}

	pp.Printf("%v\n", blkInfo)
}
