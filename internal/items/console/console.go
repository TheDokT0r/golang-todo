package console

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// Clears the terminal
func Clear() {
	switch runtime.GOOS {
	case "windows":
		log.Fatal("This program doesn't currently work on Windows ATM")
	default:
		c := exec.Command("sh", "-c", "clear")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	}
}
