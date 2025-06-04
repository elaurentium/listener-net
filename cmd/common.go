package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

)

var (
	Logger = func () *log.Logger {
		return log.New(os.Stderr, "", log.LstdFlags)
	}()

	IP_DEFAULT = ""
)

func CheckOS() {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		err := os.Chmod("./build.sh", 0755)

		if err != nil {
			fmt.Println(err)
		}

		command := exec.Command("./build.sh")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}
}