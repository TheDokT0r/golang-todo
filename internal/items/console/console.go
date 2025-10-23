package system

import (
	"log"
	"os"
	"os/exec"
	"runtime"
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
