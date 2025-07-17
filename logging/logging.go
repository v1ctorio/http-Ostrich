package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/gookit/color"
)

var debugLevel bool = false
var debug *log.Logger

func LogAndTerminate(message string, v ...any) {
	color.Yellowln(fmt.Sprintf(message, v...))
	os.Exit(1)
}

func ErrorAndKill(message string, err error) {

	println(color.Red.Render(fmt.Sprint(err))) //fun fact: the builtin print function prints to stderr, idk why

	println(color.Yellow.Renderln(message))
	os.Exit(1)
}

func SetLogLevel(setDebug bool) {

	if setDebug {
		debugLevel = true
		debug = log.New(os.Stderr, "DEBUG> ", log.Ltime)

	}
}

func DebugLog(message string, v ...any) {
	if !debugLevel {
		return
	}
	debug.Println(color.Green.Render(fmt.Sprintf(message, v...)))
}

func PrintInfoBox(address string, passphrase string, doZip bool, fLen int) {
	Box := box.New(box.Config{Px: 2, Py: 1, Type: "Round", Color: "Green", TitlePos: "Top"})

	url := address

	if passphrase != "" {
		url = fmt.Sprintf("http://u:%s@%s", passphrase, address)
	} else {
		url = fmt.Sprintf("http://%s", address)
	}
	content := `Listening on %s`
	content = fmt.Sprintf(content, url)
	Box.Print("http-Ostrich", content)
	_ = Box
}

func WarnBox(message string) {
	Box := box.New(box.Config{Px: 2, Py: 1, Type: "Bold", Color: "Yellow", TitlePos: "Inside"})

	Box.Print(color.FgYellow.Render("WARNING"), message)
}
