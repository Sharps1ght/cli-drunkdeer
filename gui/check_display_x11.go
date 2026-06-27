//go:build !wayland

package gui

import (
	"fmt"
	"os"
)

func checkDisplay() {
	if os.Getenv("DISPLAY") == "" {
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			fmt.Fprintln(os.Stderr, "Wayland session detected. Rebuild with Wayland support:\n  GO_TAGS=wayland make\nOr run with: export DISPLAY=:0 && ./drunkdeer gui")
		} else {
			fmt.Fprintln(os.Stderr, "No display server found")
		}
		os.Exit(1)
	}
}
