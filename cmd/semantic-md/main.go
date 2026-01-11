package main

import (
	"os"

	"github.com/thorstenpfister/semantic-markdown/cmd/semantic-md/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
