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
		},
		Action: func(cCtx *cli.Context) error {
			for range cCtx.Int("n") {
				fmt.Println(cuid2.CreateId())
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
