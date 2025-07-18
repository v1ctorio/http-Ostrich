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

var recursive bool

const VERSION = "0.9.1"

func main() {

	app := &cli.Command{
		Name:        "http-ostrich",
		Description: "The easy and fast http file sharing ostrich.",
		Usage:       "The http file-sharing ostrich",
		Version:     VERSION,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "Port to listen on",
				Value:   0,
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
				Usage:   "Wether to compress the files into a zip file",
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

			args := cmd.StringArgs("paths")

			if len(args) > 1 && recursive {
				logging.WarnBox("More than one argument provided but recursive flag is set. Only the first argument will be handled.")

			}

			if passphrase == "" && expose {
				logging.WarnBox("Server exposed to all incoming connections without authentication set. \nUse the --passphrase flag to setup authentication")
			}

			FilesInfo, Files, err := filemanagment.HandleFiles(args, recursive, &ShareName)
			if err != nil {
				logging.ErrorAndKill("Error reading the provided files", err)
				return nil
			}

			if doZip {
				FilesInfo, Files = filemanagment.ZipFiles(FilesInfo, Files, ShareName)
			}

			if len(Files) == 0 {
				logging.LogAndTerminate("No files found to serve. Aborting")
			}

			/*go*/

			address := web.GenerateListenAddress(port, expose)

			logging.PrintInfoBox(address, passphrase, doZip, len(FilesInfo))

			web.HttpServer(address, passphrase, Files, FilesInfo, ShareName)

			return nil

		},
		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "paths",
				Min:  0,
				Max:  -1,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
