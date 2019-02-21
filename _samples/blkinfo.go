package main

import (
	"fmt"
	"os"
	"strings"

	blkinfo "github.com/blp1526/go-blkinfo"
	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		desc := []string{
			"USAGE:",
			"  go run _samples/blkinfo.go [device]",
			"EXAMPLE:",
			"  go run _samples/blkinfo.go /dev/sda1",
		}
		fmt.Println(strings.Join(desc, "\n"))
		os.Exit(0)
	}

	path := args[1]
	bi, err := blkinfo.New(path)
	if err != nil {
		panic(err)
	}

	b, err := yaml.Marshal(bi)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
}
