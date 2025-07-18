package web

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/v1ctorio/http-ostrich/logging"
)

const DEFAULT_PORT = 8069

var ShareName string = "Ostrich shared files"

var FilesInfo []os.FileInfo
var Files []*os.File

//go:embed templates/root.html.tmpl
var rootTemplate string

type TemplateData struct {
	Title string // the shareName
	Files []os.FileInfo
}

func GenerateListenAddress(port int, expose bool) string {

	if port == 0 {

		if isPortFree(DEFAULT_PORT) {
			port = DEFAULT_PORT
		} else {

			var err error
			port, err = getFreePort()
			if err != nil {
				logging.ErrorAndKill("Error trying to get a free port", err)
			}
		}
	}

	address := fmt.Sprintf(":%d", port)
	if !expose {
		address = fmt.Sprintf("%s%s", getLocalIP(), address)
	} else {
		address = fmt.Sprintf("0.0.0.0%s", address)
	}
	return address
}

func HttpServer(address string, passphrase string, files []*os.File, filesInfo []os.FileInfo, shareName string) string {

	Files = files
	FilesInfo = filesInfo
	ShareName = shareName

	mux := http.NewServeMux()

	if passphrase != "" {

		mux.Handle("/", authMiddleware(http.HandlerFunc(getRoot), passphrase))
		mux.Handle("/dl", authMiddleware(http.HandlerFunc(getdl), passphrase))
	} else {
		mux.HandleFunc("/", getRoot)
		mux.HandleFunc("/dl", getdl)
	}

	err := http.ListenAndServe(address, mux)

	if err != nil {
		logging.ErrorAndKill("Error trying to start the file server", err)
	}

	return address
}

func getRoot(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{
		Title: ShareName,
		Files: FilesInfo,
	}
	tmpl := generateRootHTMLTemplate()
	logging.DebugLog("got / request")
	tmpl.Execute(w, data)
}
func getdl(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		logging.DebugLog("Error parsing the url %v", err)
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
		// function borrowed (stolen and slightly edited to fit my need) from https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d?permalink_comment_id=4467913#gistcomment-4467913
		"pretty_fsize": func(bytes int64) string {
			f := float64(bytes)
			for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
				if math.Abs(f) < 1024.0 {
					return fmt.Sprintf("%3.1f%sB", f, unit)
				}
				f /= 1024.0
			}
			return fmt.Sprintf("%.1fYiB", f)
		},
	}
	tmpl, err := template.New("Root").Funcs(funcMap).Parse(rootTemplate)
	if err != nil {
		logging.LogAndTerminate("Error parsing the HTML template %v", err)
	}
	return *tmpl

}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func isPortFree(port int) bool {
	if a, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port)); err == nil {
		a.Close()
		return true

	}

	return false
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		logging.ErrorAndKill("Error trying to get the listening IP address", err)
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
