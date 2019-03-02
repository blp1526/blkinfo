package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blp1526/blkinfo"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

const exitCodeNG = 1

var version string

func main() {
	version = "0.0.1"

	app := cli.NewApp()
	app.Name = "blkinfo"
	app.Usage = ""
	app.Version = version
	app.Authors = []cli.Author{
		{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 Shingo Kawamura"

	allowedFormat := "[json|yaml]"
	var showCommand = cli.Command{
		Name:      "show",
		Usage:     "Show a block device info.",
		ArgsUsage: "[path]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "format",
				Value: "json",
				Usage: fmt.Sprintf("output format %s", allowedFormat),
			},
		},
		Action: func(c *cli.Context) (err error) {
			path := c.Args().First()
			bi, err := blkinfo.New(path)
			if err != nil {
				return cli.NewExitError(err, exitCodeNG)
			}

			var b []byte
			format := c.String("format")
			switch format {
			case "json":
				b, err = json.MarshalIndent(bi, "", "  ")
				if err != nil {
					return cli.NewExitError(err, exitCodeNG)
				}
			case "yaml":
				b, err = yaml.Marshal(bi)
				if err != nil {
					return cli.NewExitError(err, exitCodeNG)
				}
			default:
				err = fmt.Errorf("unknown format '%s', expected %s", format, allowedFormat)
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
