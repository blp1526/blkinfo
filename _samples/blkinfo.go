package main

import (
	"fmt"

	blkinfo "github.com/blp1526/go-blkinfo"
	"gopkg.in/yaml.v2"
)

func main() {
	bi, err := blkinfo.New("/dev/sda1")
	if err != nil {
		panic(err)
	}

	b, err := yaml.Marshal(bi)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
}
