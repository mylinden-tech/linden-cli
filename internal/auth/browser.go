package auth

import (
	"os/exec"
	"runtime"
)

// openBrowser opens the given URL in the default browser.
func openBrowser(rawURL string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{rawURL}
	case "linux":
		cmd = "xdg-open"
		args = []string{rawURL}
	default:
		return
	}

	_ = exec.Command(cmd, args...).Start()
}
