package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	var port int
	var expose bool
	doZip := false
	passphrase := ""
	var files []os.File

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

			fileInfo, err := os.Stat(destinationPath)

			if err != nil {
				log.Fatalf("Error getting file info: %v", err)
			}

			isDirectory := fileInfo.IsDir()

			files = handleFiles(fileInfo)

			if isDirectory {

				if zip {
					return nil
				}

				dir, err := os.Open(destinationPath)
				if err != nil {
					log.Fatalf("Error opening directory: %v", err)
				}
				files, err := dir.ReadDir(0)
				if err != nil {
					log.Fatalf("Error reading directory: %v", err)
				}

				for _, file := range files {

				}

			} else {
				file, err := os.Open(destinationPath)
				if err != nil {
					log.Fatalf("Error opening file: %v", err)
				}
				defer file.Close()

				files = append(files, *file)
			}

			fmt.Printf("Starting server on port %d\n serving a directory", port, isDirectory)
			// Here you would start your server
			return nil

		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func handleFiles(path os.FileInfo) ([]os.File, error) {

	if os.FileInfo.IsDir(path) {
		if doZip {

		}

	}

	return nil, nil

}
