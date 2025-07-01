package logging

import (
	"fmt"
	"log"
	"os"
)

var debugLevel bool = false
var debug *log.Logger
var NoColor = os.Getenv("NO_COLOR") == ""

const END = "\033[0m"
const GREEN = "\033[32m"
const RED = "\033[31m"
const YELLOW = "\033[33m"

func LogAndTerminate(message string, v ...any) {
	println(fmt.Sprintf(message, v...))
	print("\n")
	os.Exit(1)
}

func ErrorAndKill(message string, err error) {
	if !NoColor {
		print("\033[1m")
	}
	println(err, "\n", "\n") //fun fact: the builtin print function prints to stderr, idk why

	if !NoColor {
		print("\033[32m" + "\n")
		print(YELLOW)
	}
	println(message, END)
	os.Exit(1)
}

func SetLogLevel(setDebug bool) {

	if setDebug {
		debugLevel = true
		if NoColor {
			debug = log.New(os.Stderr, GREEN+"", log.Ltime)

		} else {
			debug = log.New(os.Stderr, "DEBUG> ", log.Ltime)
		}
	}
}

func DebugLog(message string, v ...any) {
	if !debugLevel {
		return
	}
	debug.Printf(message, v...)
	debug.Print("\n", END)
}
