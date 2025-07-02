package main

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	filemanagment "github.com/v1ctorio/http-ostrich/file-managment"
	"github.com/v1ctorio/http-ostrich/logging"
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
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "",
				Value:   false,
				Aliases: []string{"v"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			doZip := cmd.Bool("zip")
			port := cmd.Int("port")
			expose := cmd.Bool("expose")
			ShareName := "Shared files"

			logging.SetLogLevel(cmd.Bool("verbose"))

			passphrase := cmd.String("passphrase")
			if cmd.Args().Len() < 1 {
				logging.LogAndTerminate("No files provided")
			}
			if cmd.Args().Len() > 1 && recursive {
				println("More than one argument provided but recursive flag is set. Only the first argument will be handled.")
			}

			args := cmd.Args().Slice()

			FilesInfo, Files, err := filemanagment.HandleFiles(args, recursive, &ShareName)
			if err != nil {
				log.Fatalf("%v", err)
				return nil
			}

			if doZip {
				FilesInfo, Files = filemanagment.ZipFiles(FilesInfo, Files)
			}

			if len(Files) == 0 {
				logging.LogAndTerminate("No files found to serve. Aborting")
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
