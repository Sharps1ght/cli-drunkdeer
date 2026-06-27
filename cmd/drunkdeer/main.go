package main

import (
	"os"

	"github.com/2xxn/cli-drunkdeer/internal/drunkdeer"
	"github.com/2xxn/cli-drunkdeer/gui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "gui" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		gui.Run()
	} else {
		drunkdeer.Run()
	}
}
