package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

type APIFile struct {
	size int
	name int //URL encoded name of the file
}

var shareName string = "Shared files"

var files []os.File
var filesInfo []os.FileInfo

//go:embed templates/root.html.tmpl
var rootTemplate string

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

			filesish, err := handleFiles(destinationPath)
			if err != nil {
				log.Fatalf("%v", err)
				return nil
			}
			files = filesish

			/*go*/
			httpServer(port, expose)
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
		log.Fatal("Error expanding")
	}

	println("%v expanded to, %v", path, fullPath)
	fileInfo, err := os.Stat(fullPath)

	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	//TODO: handle zipping
	isDirectory := fileInfo.IsDir()

	if isDirectory {

		filesList, err := os.ReadDir(fullPath)
		if err != nil {
			return nil, err
		}
		var filesToReturn []os.File

		err = os.Chdir(fullPath)

		if err != nil {
			log.Fatal("Error opening the directory", fullPath)
		}

		for i, e := range filesList {
			fInfo, err := e.Info()

			filesInfo = append(filesInfo, fInfo)
			if err != nil {
				log.Fatal(err)
			}
			if fInfo.IsDir() {
				continue
			}
			file, err := os.OpenFile(fInfo.Name(), os.O_RDONLY, e.Type().Perm())
			if err != nil {
				log.Fatal(err)
			}
			filesToReturn = append(filesToReturn, *file)
			println(i, e)

		}

		shareName = fileInfo.Name()
		return filesToReturn, nil

	}
	file, err := os.OpenFile(fullPath, os.O_RDONLY, fileInfo.Mode().Perm())

	uniqueFileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	filesInfo = []os.FileInfo{uniqueFileInfo}

	if err != nil {
		return nil, err
	}

	shareName = fileInfo.Name()
	return []os.File{*file}, nil

}

func httpServer(port int, expose bool) {

	//TODO: implement not-expose
	address := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	//mux.HandleFunc("/dl", getdl)

	println("Started listening in ", address)
	err := http.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}
}

type TemplateData struct {
	Title string // the shareName
	Files []os.FileInfo
}

func getRoot(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{
		Title: shareName,
		Files: filesInfo,
	}
	tmpl := generateRootHTMLTemplate()
	fmt.Print("got / request\n", data.Title)
	tmpl.Execute(w, data)
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func generateRootHTMLTemplate() template.Template {
	tmpl, err := template.New("Root").Parse(rootTemplate)
	if err != nil {
		log.Fatalf("Error parsing the HTML template %d", err)
	}
	return *tmpl

}
