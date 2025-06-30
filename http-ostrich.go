package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	filemanagment "github.com/v1ctorio/http-ostrich/file-managment"
	"github.com/v1ctorio/http-ostrich/web"
)

type APIFile struct {
	size int
	name int //URL encoded name of the file
}

var recursive bool

func main() {

	app := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "Port to listen on",
				Value:   8069,
				Aliases: []string{"p"},
			},
			&cli.BoolFlag{
				Name:    "expose",
				Usage:   "Wether to expose the server to foreign IPs",
				Value:   false,
				Aliases: []string{"e"},
			},
			&cli.StringFlag{
				Name:    "passphrase",
				Usage:   "Passphrase for basic authentication",
				Value:   "",
				Aliases: []string{"a"},
			},
			&cli.BoolFlag{
				Name:    "zip",
				Usage:   "Wether to zip the files",
				Value:   false,
				Aliases: []string{"z"},
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Usage:       "",
				Value:       false,
				Aliases:     []string{"r"},
				Destination: &recursive,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			doZip := cmd.Bool("zip")
			_ = doZip
			port := cmd.Int("port")
			expose := cmd.Bool("expose")
			ShareName := "Shared files"

			passphrase := cmd.String("passphrase")
			if cmd.Args().Len() < 1 {
				fmt.Println("No command provided, exiting")
				return nil
			}

			args := cmd.Args().Slice()

			FilesInfo, Files, err := filemanagment.HandleFiles(args, recursive, &ShareName)
			if err != nil {
				log.Fatalf("%v", err)
				return nil
			}

			if len(Files) == 0 {
				log.Fatalf("No files found to server. Aborting")
			}

			/*go*/

			web.HttpServer(port, expose, passphrase, Files, FilesInfo, ShareName)
			return nil

		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func logBox() {
	// TODO: fancy cool looking box with the stuff and that
}
