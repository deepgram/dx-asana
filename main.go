package main

import (
	"os"

	"github.com/deepgram/dx-asana/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}