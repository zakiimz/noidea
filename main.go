// main.go - Entry point for the noidea application
//
// noidea is a Git companion that provides commit message suggestions,
// feedback on your commit history, and other helpful development tools.
//
// This main package initializes the application and invokes the root command.

package main

import (
	"github.com/AccursedGalaxy/noidea/cmd"
)

func main() {
	cmd.Execute()
}
