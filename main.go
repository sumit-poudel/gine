package main

import (
	"fmt"
	"os"

	"github.com/sumit-poudel/gine/internal/app"
)

func main() {
	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "gine: %v\n", err)
		os.Exit(1)
	}
}
