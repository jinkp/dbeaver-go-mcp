package dbeaver

import (
	"os/exec"
	"runtime"
	"strings"
)

// IsRunning returns true if a DBeaver process appears to be running.
// On any error executing the check command, returns false (fail open —
// don't block installation if the process check itself fails).
func IsRunning() bool {
	switch runtime.GOOS {
	case "windows":
		names := []string{"dbeaver.exe", "dbeaver64.exe", "dbeaver-ce.exe"}
		for _, name := range names {
			out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq "+name).Output()
			if err == nil && strings.Contains(string(out), name) {
				return true
			}
		}
		return false
	default:
		return exec.Command("pgrep", "-x", "dbeaver").Run() == nil
	}
}
