package main

import (
	"fmt"

	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/cmd/sub"
)



func init() {
	cmd.CheckOS()
	cmd.PrintBanner()
}

func main() {
	fmt.Printf("Welcome to the Listener Net application!\n")
	sub.Interfaces()
}