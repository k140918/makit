package main

import (
	"github.com/k140918/makit/cmd"
)

func main() {
	cmd.Version = VERSION
	cmd.Execute()
}
