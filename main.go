package main

import "github.com/jheddings/ccglow/cmd"

func main() {
	cmd.SetVersion(cmd.Version())
	cmd.Execute()
}
