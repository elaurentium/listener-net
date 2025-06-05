package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

)

var (
	Logger = func () *log.Logger {
		return log.New(os.Stderr, "", log.LstdFlags)
	}()

	IP_DEFAULT = ""
)

func CheckOS() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}
}