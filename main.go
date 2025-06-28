package main

import (
	"context"
	_ "embed"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

type APIFile struct {
	size int
	name int //URL encoded name of the file
}

var shareName string = "Shared files"

var filesInfo []os.FileInfo
var files []*os.File

//go:embed templates/root.html.tmpl
var rootTemplate string

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

			if len(files) == 0 {
				log.Fatalf("No files found to server. Aborting")
			}

			/*go*/

			httpServer(port, expose, passphrase)
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
			file, err := os.OpenFile(e.Name(), os.O_RDONLY, fInfo.Mode().Perm())
			panic(err)
			files = append(files, file)

			filesInfo = append(filesInfo, fInfo)
			panic(err)
			if fInfo.IsDir() {
				continue
			}
			println(i, e)

		}

		shareName = dirInfo.Name()
		return nil
	}

	for _, fileName := range args {

		fullPath, err := filepath.Abs(fileName)
		fmt.Println("Expanded ", fileName, fullPath)
		if err != nil {
			return errors.New("error expanding the provided path")
		}

		file, err := os.OpenFile(fullPath, os.O_RDONLY, 0)
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

		filesInfo = append(filesInfo, fInfo)
		files = append(files, file)
	}

	return nil

}

func httpServer(port int, expose bool, passphrase string) {

	//TODO: implement not-expose
	address := fmt.Sprintf(":%d", port)

	if !expose {
		address = fmt.Sprintf("localhost%s", address)
	}
	mux := http.NewServeMux()

	if passphrase != "" {

		mux.Handle("/", authMiddleware(http.HandlerFunc(getRoot), passphrase))
		mux.Handle("/dl", authMiddleware(http.HandlerFunc(getdl), passphrase))
	} else {
		mux.HandleFunc("/", getRoot)
		mux.HandleFunc("/dl", getdl)
	}

	println("Started listening in ", address, " with authentication: ", passphrase)
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
func getdl(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print(err)
		return
	}

	params, _ := url.ParseQuery(parsedURL.RawQuery)

	requestedFile := params.Get("f")

	if requestedFile == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "File to be downloaded not specified")
		return
	}

	requestedFile, _ = url.QueryUnescape(requestedFile)

	var allowedFile int = -1
	for i, f := range filesInfo {
		if f.Name() == requestedFile {
			allowedFile = i
		}
	}
	if allowedFile == -1 {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "File not found or not shared")
		return
	}

	allowedFileInfo := filesInfo[allowedFile]
	file := files[allowedFile]

	defer file.Close()

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", allowedFileInfo.Name()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)

	fmt.Printf("got /dl request for %s\n", requestedFile)
}

func generateRootHTMLTemplate() template.Template {

	funcMap := template.FuncMap{
		"escape": func(s string) string {
			return url.QueryEscape(s)
		},
	}
	tmpl, err := template.New("Root").Funcs(funcMap).Parse(rootTemplate)
	if err != nil {
		log.Fatalf("Error parsing the HTML template %d", err)
	}
	return *tmpl

}

func authMiddleware(next http.Handler, passphrase string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"401\"")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Authentication required.")
			return
		}
		println("Unparse auth ", auth)
		auth = strings.Split(auth, " ")[1] // in "Basic <encoded>" only keep the encoded part
		dec, err := b64.StdEncoding.DecodeString(auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Error parsing the authentication")
			return
		}
		auth = string(dec)
		auth = strings.Split(auth, ":")[1] // in "user:pass" only save the pass

		fmt.Printf("Tried to access using %s", auth)
		if auth != passphrase {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"401\"")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Authentication required.")
			return
		}

		next.ServeHTTP(w, r)

	})
}

func logBox() {
}

func panic(err error) {

	if err != nil {
		log.Fatalf("Panic! %d \n", err)
	}
}
