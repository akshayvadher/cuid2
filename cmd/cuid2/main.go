package main

import (
	"fmt"
	"github.com/akshayvadher/cuid2"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "cuid2",
		Usage: "Create CUID2",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "n",
				Value: 1,
				Usage: "Numbers of ids to generate",
			},
			&cli.IntFlag{
				Name:  "len",
				Value: 24,
				Usage: "Length of the Id (between 2 and 36)",
				Action: func(ctx *cli.Context, v int) error {
					if v > 36 || v < 2 {
						return fmt.Errorf("len %v should be between 2 and 36", v)
					}
					return nil
				},
			},
		},
		Action: func(cCtx *cli.Context) error {
			for range cCtx.Int("n") {
				fmt.Println(cuid2.CreateIdOf(cCtx.Int("len")))
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
