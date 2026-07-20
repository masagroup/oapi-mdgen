package main

import (
	"context"
	"log"
	"os"

	"github.com/masagroup/oapi-mdgen/driver"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "oapi-mdgen",
		Usage: "Convert OpenAPI specifications to Markdown documentation",
		Description: "The command lets you convert OpenAPI specifications to markdown documentation.\n" +
			"More information on command can be found here: https://github.com/masagroup/oapi-mdgen",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "openapi input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "markdown output file",
				Required: true,
			}},
		UseShortOptionHandling: true,
		Action: func(ctx context.Context, c *cli.Command) error {
			input := c.String("input")
			output := c.String("output")
			if input == "" || output == "" {
				return cli.Exit("input and output flags are required", 1)
			}
			return driver.GenerateMarkdown(input, output)
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
