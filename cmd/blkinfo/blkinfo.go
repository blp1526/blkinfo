package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blp1526/blkinfo"
	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
)

const exitCodeNG = 1

var version string

func main() {
	version = "0.0.1"

	app := cli.NewApp()
	app.Name = "blkinfo"
	app.Usage = "block device information utility for Linux"
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
		Usage:     "Shows block device information",
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

			b, err := json.MarshalIndent(bi, "", "  ")
			if err != nil {
				return cli.NewExitError(err, exitCodeNG)
			}

			format := c.String("format")
			switch format {
			case "json":
				break
			case "yaml":
				b, err = yaml.JSONToYAML(b)
				if err != nil {
					return cli.NewExitError(err, exitCodeNG)
				}
			default:
				err = fmt.Errorf("unknown format '%s', expected %s", format, allowedFormat)
				return cli.NewExitError(err, exitCodeNG)
			}

			s := strings.TrimSpace(string(b))
			fmt.Printf("%s\n", s)
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
