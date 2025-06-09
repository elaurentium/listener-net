package main

import (
	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/cmd/sub"
)



func init() {
	cmd.CheckOS()
	cmd.PrintBanner()
}

func main() {
	cmd.Usage("Welcome to the Listener Net application!\n")
	sub.Interfaces(nil)
}