//go:build wayland

package gui

import (
	"fmt"
	"os"
)

func checkDisplay() {
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		fmt.Fprintln(os.Stderr, "Wayland build requires a Wayland compositor")
		os.Exit(1)
	}
}
