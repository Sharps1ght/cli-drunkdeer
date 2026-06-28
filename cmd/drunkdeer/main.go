package main

import (
	"os"
	"runtime"

	"github.com/2xxn/cli-drunkdeer/internal/drunkdeer"
	"github.com/2xxn/cli-drunkdeer/gui"
)

func main() {
	launchGUI := false
	if len(os.Args) > 1 && os.Args[1] == "gui" {
		launchGUI = true
		os.Args = append(os.Args[:1], os.Args[2:]...)
	} else if len(os.Args) == 1 && runtime.GOOS == "windows" {
		launchGUI = true
	}

	if launchGUI {
		gui.Run()
	} else {
		drunkdeer.Run()
	}
}
