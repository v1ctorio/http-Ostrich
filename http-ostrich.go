package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
	"v1c.rocks/http-ostrich/web"
)

type APIFile struct {
	size int
	name int //URL encoded name of the file
}

var ShareName string = "Shared files"

var FilesInfo []os.FileInfo
var Files []*os.File

var doZip bool = false
var recursive bool = false

func main() {

	var port int
	var expose bool

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
				Name:    "passphrase",
				Usage:   "Passphrase for basic authentication",
				Value:   "",
				Aliases: []string{"a"},
			},
			&cli.BoolFlag{
				Name:        "zip",
				Usage:       "Wether to zip the files",
				Value:       false,
				Aliases:     []string{"z"},
				Destination: &doZip,
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

			passphrase := cmd.String("passphrase")
			if cmd.Args().Len() < 1 {
				fmt.Println("No command provided, exiting")
				return nil
			}

			args := cmd.Args().Slice()

			err := handleFiles(args)
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

func handleFiles(args []string) error {

	if recursive {
		directory := args[0]
		fullPath, err := filepath.Abs(directory)
		if err != nil {
			log.Fatal("Error expanding")
		}
		println("%v expanded to, %v", directory, fullPath)
		dirInfo, err := os.Stat(fullPath)
		if err != nil {
			log.Fatalf("Error getting dir info: %v", err)
		}

		//TODO: handle zipping
		isDirectory := dirInfo.IsDir()

		if !isDirectory {
			return errors.New("directory was not provided but recursive flag is set.")
		}
		filesList, err := os.ReadDir(fullPath)
		if err != nil {
			return err
		}
		err = os.Chdir(fullPath)

		if err != nil {
			log.Fatal("Error opening the directory", fullPath)
		}

		for i, e := range filesList {
			fInfo, err := e.Info()
			panic(err)
			file, err := os.Open(e.Name())
			panic(err)
			Files = append(Files, file)

			FilesInfo = append(FilesInfo, fInfo)
			panic(err)
			if fInfo.IsDir() {
				continue
			}
			println(i, e)

		}

		ShareName = dirInfo.Name()
		return nil
	}

	for _, fileName := range args {

		fullPath, err := filepath.Abs(fileName)
		fmt.Println("Expanded ", fileName, fullPath)
		if err != nil {
			return errors.New("error expanding the provided path")
		}

		file, err := os.Open(fullPath)
		if err != nil {
			return err
		}

		fInfo, err := file.Stat()
		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			log.Println("Skipping file because it is a directory and the recursive flag is not set ", file.Name())
			file.Close()
			continue
		}

		FilesInfo = append(FilesInfo, fInfo)
		Files = append(Files, file)
	}

	return nil

}

type TemplateData struct {
	Title string // the shareName
	Files []os.FileInfo
}

func logBox() {
	// TODO: fancy cool looking box with the stuff and that
}

func panic(err error) {

	if err != nil {
		log.Fatalf("Panic! %d \n", err)
	}
}
