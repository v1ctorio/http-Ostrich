package web

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var ShareName string = "Shared files"

var FilesInfo []os.FileInfo
var Files []*os.File

//go:embed templates/root.html.tmpl
var rootTemplate string

type TemplateData struct {
	Title string // the shareName
	Files []os.FileInfo
}

func HttpServer(port int, expose bool, passphrase string, files []*os.File, filesInfo []os.FileInfo, shareName string) {

	Files = files
	FilesInfo = filesInfo
	ShareName = shareName

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

func getRoot(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{
		Title: ShareName,
		Files: FilesInfo,
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
	for i, f := range FilesInfo {
		if f.Name() == requestedFile {
			allowedFile = i
		}
	}
	if allowedFile == -1 {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "File not found or not shared")
		return
	}

	allowedFileInfo := FilesInfo[allowedFile]
	file := Files[allowedFile]

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
