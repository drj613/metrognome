package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drj613/metrognome/internal/ui"
)

func main() {
	fmt.Println("🎩 Welcome to Metrognome - Where Every Beat is Garden Fresh! 🌱")
	fmt.Println()

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start the garden metronome: %v", err)
	}
}
