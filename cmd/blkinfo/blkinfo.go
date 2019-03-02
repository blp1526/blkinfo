package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blp1526/blkinfo"
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
		Usage:     "Show a block device info.",
		ArgsUsage: "[device]",
		Action: func(c *cli.Context) (err error) {
			path := c.Args().First()
			bi, err := blkinfo.New(path)
			if err != nil {
				return cli.NewExitError(err, exitCodeNG)
			}

			b, err := json.MarshalIndent(bi, "", "  ")
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
