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

func main() { // nolint: funlen
	app := cli.NewApp()
	app.Name = "blkinfo"
	app.Usage = "block device information utility for Linux"
	app.UsageText = "blkinfo [options] path"
	app.Version = blkinfo.Version()
	app.Description = fmt.Sprintf("commit %s, built at %s", blkinfo.Revision(), blkinfo.BuiltAt())
	app.Copyright = "(c) 2019 Shingo Kawamura"
	app.Authors = []cli.Author{
		{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}
	app.HideHelp = true
	allowedFormat := "[json|yaml]"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: fmt.Sprintf("output format %s", allowedFormat),
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}

		if len(c.Args()) != 1 {
			cli.ShowAppHelpAndExit(c, 0)
		}

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
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
