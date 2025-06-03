package main

import (
	"fmt"

	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/cmd/sub"
)

func init() {
	cmd.PrintBanner()
	sub.Interfaces()
}

func main() {
	fmt.Printf("Welcome to the Listener Net application!\n")
}