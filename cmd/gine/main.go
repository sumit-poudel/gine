package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/sumit-poudel/gine/internal/tui"
)

func main() {
	if _, err := tea.NewProgram(tui.New()).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
