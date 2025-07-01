package web

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/v1ctorio/http-ostrich/logging"
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
		logging.ErrorAndKill("Error trying to start the file server", err)
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
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error parsind the url %v", err))
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

	logging.DebugLog("got /dl request for %s", requestedFile)
}
func generateRootHTMLTemplate() template.Template {

	funcMap := template.FuncMap{
		"escape": func(s string) string {
			return url.QueryEscape(s)
		},
	}
	tmpl, err := template.New("Root").Funcs(funcMap).Parse(rootTemplate)
	if err != nil {
		logging.LogAndTerminate("Error parsing the HTML template %v", err)
	}
	return *tmpl

}
