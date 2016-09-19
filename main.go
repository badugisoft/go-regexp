package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/badugisoft/xson"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := cli.App{
		Name:    "go-regexp",
		Usage:   "go regular expression commandline tool",
		Version: "0.0.1",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Value:   "yaml",
				Usage:   "output format",
			},
		},
		ArgsUsage: "[name=regexp]",
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return cli.Exit("at least one regexp is required", -1)
	}

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return cli.Exit("reading input failed", -2)
	}

	outputType := xson.GetType(c.String("type"))
	if outputType == xson.UNKNOWN {
		return cli.Exit("unsupported output format: "+c.String("type"), -3)
	}

	text := string(in)
	res := map[string][][]string{}
	for _, arg := range c.Args().Slice() {
		pos := strings.Index(arg, "=")
		if pos < 0 {
			cli.Exit("invalid argument: "+arg, -4)
		}

		r, err := regexp.Compile(arg[pos+1:])
		if err != nil {
			return cli.Exit("invalid regexp: "+arg, -5)
		}

		res[arg[0:pos]] = r.FindAllStringSubmatch(text, -1)
	}

	out, err := xson.Marshal(outputType, res)

	fmt.Println(string(out))

	return nil
}
