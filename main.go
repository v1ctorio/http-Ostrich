package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func main() {

	var port int
	var expose bool
	doZip := false
	passphrase := ""

	app := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Port to listen on",
				Value:       8069,
				Destination: &port,
				Aliases:     []string{"p"},
			},
			&cli.BoolFlag{
				Name:        "expose",
				Usage:       "Wether to expose the server to foreign IPs",
				Value:       false,
				Destination: &expose,
				Aliases:     []string{"e"},
			},
			&cli.StringFlag{
				Name:        "passphrase",
				Usage:       "Passphrase for basic authentication",
				Value:       "",
				Destination: &passphrase,
				Aliases:     []string{"a"},
			},
			&cli.BoolFlag{
				Name:        "zip",
				Usage:       "Wether to zip the files",
				Value:       false,
				Aliases:     []string{"z"},
				Destination: &doZip,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			if cmd.Args().Len() < 1 {
				fmt.Println("No command provided, exiting")
				return nil
			}

			destinationPath := cmd.Args().First()

			files, err := handleFiles(destinationPath)
			_ = files

			if err != nil {
				log.Fatalf("%v", err)
				return nil
			}

			fmt.Printf("Starting server on port %d\n serving a directory", port)
			// Here you would start your server
			return nil

		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func handleFiles(path string) ([]os.File, error) {

	fullPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}
	println("%v expanded to, %v", path, fullPath)
	fileInfo, err := os.Stat(fullPath)

	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	isDirectory := fileInfo.IsDir()

	_ = isDirectory

	return nil, nil

}
