[![Build Status](https://travis-ci.org/blp1526/blkinfo.svg?branch=master)](https://travis-ci.org/blp1526/blkinfo)

# blkinfo

A Linux Block Device Info Library

## Usage

### As A Package

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/blp1526/blkinfo"
)

func main() {
	bi, _ := blkinfo.New("/dev/sda")
	b, _ := json.Marshal(bi)
	fmt.Println(string(b))
}
```

### As A CLI

```
$ make
$ ./bin/blkinfo show /dev/sda
```
