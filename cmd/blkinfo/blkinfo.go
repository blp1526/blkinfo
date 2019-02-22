package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blp1526/go-blkinfo"
	"github.com/urfave/cli"
)

const exitCodeNG = 1

var version string
var revision string

func main() {
	app := cli.NewApp()
	app.Name = "blkinfo"
	app.Usage = ""
	app.Version = version
	app.Description = fmt.Sprintf("REVISION: %s", revision)
	app.Authors = []cli.Author{
		{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 Shingo Kawamura"

	var showCommand = cli.Command{
		Name:      "show",
		Usage:     "show blkinfo",
		ArgsUsage: "path",
		Action: func(c *cli.Context) (err error) {
			path := c.Args().First()
			bi, err := blkinfo.New(path)
			if err != nil {
				return cli.NewExitError(err, exitCodeNG)
			}

			b, err := json.Marshal(bi)
			if err != nil {
				return cli.NewExitError(err, exitCodeNG)
			}

			fmt.Printf("%s\n", string(b))

			return nil
		},
	}

	app.Commands = []cli.Command{
		showCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
