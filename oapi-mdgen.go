package main

import (
	"context"
	"fmt"
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
			inputFlag := c.String("input")
			outputFlag := c.String("output")
			if inputFlag == "" || outputFlag == "" {
				return cli.Exit("input and output flags are required", 1)
			}

			// load an OpenAPI 3 specification from bytes
			input, err := os.Open(inputFlag)
			if err != nil {
				return fmt.Errorf("cannot read input file: %w", err)
			}

			// create output file
			output, err := os.Create(outputFlag)
			if err != nil {
				return fmt.Errorf("cannot create output file: %w", err)
			}
			defer func() {
				if err2 := output.Close(); err2 != nil && err == nil {
					err = fmt.Errorf("cannot close output file: %w", err2)
				}
			}()

			return driver.GenerateMarkdown(input, output)
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
