package main

import (
	"os"

	"github.com/linzeyan/gpgen/cmd"
)

func main() {
	if err := cmd.Run().Execute(); err != nil {
		os.Exit(1)
	}
}
