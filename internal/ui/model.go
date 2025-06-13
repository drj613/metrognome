package ui

import (
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drj613/metrognome/internal/metronome"
)

// Model represents the UI state
type Model struct {
	metronome      *metronome.Metronome
	currentBeat    int
	lastBeatTime   time.Time
	selectedPreset int
	showPresets    bool
	showHelp       bool
	help           help.Model
	commandsTable  table.Model
	keys           keyMap
	width          int
	height         int
	beatAnimation  int
	gnomeFrame     int
	soundEnabled   bool
	starColors     []string
	starPositions  [][]int
	pendulumAngle  float64
}

// beatMsg is sent when a beat occurs
type beatMsg int

// tickMsg is for animations
type tickMsg time.Time

// keyMap defines our key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Space  key.Binding
	Tab    key.Binding
	Preset key.Binding
	Sound  key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Space, k.Tab, k.Sound, k.Preset, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Space, k.Tab, k.Sound, k.Preset},
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "increase BPM"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "decrease BPM"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous time signature"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next time signature"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "start/stop"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle time signatures"),
	),
	Preset: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "toggle presets"),
	),
	Sound: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle sound"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// createCommandsTable creates a styled table with commands
func createCommandsTable() table.Model {
	columns := []table.Column{
		{Title: "Key", Width: 12},
		{Title: "Action", Width: 30},
		{Title: "Gnome's Wisdom", Width: 40},
	}

	rows := []table.Row{
		{"Space", "Start/Stop metronome", "Every gnome needs their rhythm!"},
		{"↑/k", "Increase BPM (+5)", "Faster steps through the garden"},
		{"↓/j", "Decrease BPM (-5)", "Slower pace for flower sniffing"},
		{"←/h", "Previous time signature", "Try different garden dances"},
		{"→/l", "Next time signature", "Explore more rhythmic patterns"},
		{"Tab", "Cycle time signatures", "Quick tempo style changes"},
		{"p", "Toggle presets menu", "Choose pre-made garden rhythms"},
		{"s", "Toggle sound on/off", "Gnomes prefer quiet sometimes"},
		{"?", "Toggle this help", "Wisdom from the garden gnome"},
		{"q/Ctrl+C", "Quit application", "Return to the mushroom house"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("86")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("86"))

	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("212")).
		Background(lipgloss.Color("236")).
		Bold(false)

	t.SetStyles(tableStyle)
	return t
}

// NewModel creates a new UI model
func NewModel() Model {
	m := Model{
		metronome:      metronome.New(120, metronome.CommonTimeSignatures[0]),
		selectedPreset: 0,
		showPresets:    false,
		showHelp:       false,
		help:           help.New(),
		commandsTable:  createCommandsTable(),
		keys:           keys,
		gnomeFrame:     0,
		soundEnabled:   true,
		starColors:     []string{"240", "244", "250", "254", "230", "226", "222", "86", "212", "231"},
	}
	m.help.ShowAll = false
	m.initializeStars()
	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenForBeats(m.metronome),
		tickAnimation(),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.commandsTable.SetWidth(msg.Width - 4)
		m.initializeStars() // Reinitialize stars when window size changes

	case beatMsg:
		m.currentBeat = int(msg)
		m.lastBeatTime = time.Now()
		m.beatAnimation = 5 // Start beat animation

		// Play sound if enabled
		if m.soundEnabled {
			// Play different sound for first beat
			isFirstBeat := m.currentBeat == 1
			go playSound(isFirstBeat)
		}

		return m, listenForBeats(m.metronome)

	case tickMsg:
		// Update animations
		if m.beatAnimation > 0 {
			m.beatAnimation--
		}
		m.gnomeFrame = (m.gnomeFrame + 1) % 4
		
		// Update pendulum swing
		if m.metronome.IsPlaying {
			// Swing based on BPM - faster BPM = faster swing
			swingSpeed := float64(m.metronome.BPM) / 60.0 * 3.14159 / 10.0
			m.pendulumAngle += swingSpeed
		}
		
		return m, tickAnimation()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.metronome.Stop()
			return m, tea.Quit

		case key.Matches(msg, m.keys.Space):
			if m.metronome.IsPlaying {
				m.metronome.Stop()
			} else {
				m.metronome.Start()
			}
			return m, listenForBeats(m.metronome)

		case key.Matches(msg, m.keys.Up):
			m.metronome.SetBPM(m.metronome.BPM + 5)
			// Reset beat animation state when BPM changes
			m.beatAnimation = 0
			m.currentBeat = 1
			// Restart beat listening if metronome was playing
			if m.metronome.IsPlaying {
				return m, listenForBeats(m.metronome)
			}

		case key.Matches(msg, m.keys.Down):
			m.metronome.SetBPM(m.metronome.BPM - 5)
			// Reset beat animation state when BPM changes
			m.beatAnimation = 0
			m.currentBeat = 1
			// Restart beat listening if metronome was playing
			if m.metronome.IsPlaying {
				return m, listenForBeats(m.metronome)
			}

		case key.Matches(msg, m.keys.Tab):
			// Cycle through time signatures
			currentIndex := 0
			for i, ts := range metronome.CommonTimeSignatures {
				if ts.Beats == m.metronome.TimeSignature.Beats &&
					ts.BeatValue == m.metronome.TimeSignature.BeatValue {
					currentIndex = i
					break
				}
			}
			nextIndex := (currentIndex + 1) % len(metronome.CommonTimeSignatures)
			m.metronome.SetTimeSignature(metronome.CommonTimeSignatures[nextIndex])
			// Reset beat animation state when time signature changes
			m.beatAnimation = 0
			m.currentBeat = 1
			// Restart beat listening if metronome was playing
			if m.metronome.IsPlaying {
				return m, listenForBeats(m.metronome)
			}

		case key.Matches(msg, m.keys.Preset):
			m.showPresets = !m.showPresets
			m.showHelp = false

		case key.Matches(msg, m.keys.Sound):
			m.soundEnabled = !m.soundEnabled

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.showPresets = false

		case key.Matches(msg, m.keys.Left):
			if m.showPresets && m.selectedPreset > 0 {
				m.selectedPreset--
			}

		case key.Matches(msg, m.keys.Right):
			if m.showPresets && m.selectedPreset < len(metronome.CommonPresets)-1 {
				m.selectedPreset++
			}

		case msg.Type == tea.KeyEnter:
			if m.showPresets {
				preset := metronome.CommonPresets[m.selectedPreset]
				m.metronome.SetBPM(preset.BPM)
				m.metronome.SetTimeSignature(preset.TimeSignature)
				m.showPresets = false
				// Reset beat animation state when preset changes
				m.beatAnimation = 0
				m.currentBeat = 1
				// Restart beat listening if metronome was playing
				if m.metronome.IsPlaying {
					return m, listenForBeats(m.metronome)
				}
			}
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	if m.showPresets {
		return m.renderPresets()
	}

	return m.renderMainWithBorder()
}

// listenForBeats creates a command that listens for metronome beats
func listenForBeats(metro *metronome.Metronome) tea.Cmd {
	return func() tea.Msg {
		beat := <-metro.BeatChannel()
		return beatMsg(beat)
	}
}

// tickAnimation creates a command for animation updates
func tickAnimation() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// renderPresets renders the preset selection view
func (m Model) renderPresets() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		MarginBottom(2)

	presetStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		MarginBottom(1)

	selectedStyle := presetStyle.Copy().
		Foreground(lipgloss.Color("212")).
		Background(lipgloss.Color("236")).
		Bold(true)

	title := titleStyle.Render("🎵 Choose Your Garden Rhythm 🎵")

	presets := ""
	for i, preset := range metronome.CommonPresets {
		style := presetStyle
		if i == m.selectedPreset {
			style = selectedStyle
		}

		line := fmt.Sprintf("%s - %d BPM (%s)\n  %s",
			preset.Name,
			preset.BPM,
			preset.TimeSignature.Name,
			preset.Description)

		presets += style.Render(line) + "\n"
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(2).
		Render("Use ←/→ to select, ENTER to confirm, P to go back")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		presets,
		instructions,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// renderHelp renders the help view
func (m Model) renderHelp() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		MarginBottom(2)

	title := titleStyle.Render("🌻 Garden Gnome's Command Guide 🌻")

	// Render the commands table
	tableView := m.commandsTable.View()

	gnomeWisdom := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginTop(2).
		Render("\"A gnome without rhythm is like a garden without flowers!\"")

	backInstruction := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		MarginTop(1).
		Render("Press '?' again to return to the garden")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		tableView,
		gnomeWisdom,
		backInstruction,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// getGnome returns a gnome for a specific beat position with alternating poses
func (m Model) getGnome(beatPosition int) string {
	// Alternate between arms down and arms up poses
	var gnome string
	if beatPosition%2 == 1 {
		// Odd positions (1, 3, 5, ...) - arms down
		gnome = "  △  \n ಠ_ಠ \n /|\\ \n / \\ "
	} else {
		// Even positions (2, 4, 6, ...) - arms up
		gnome = "  △  \n ಠ_ಠ \n \\|/ \n / \\ "
	}
	
	// Check if this gnome should be lit up for the current beat
	if m.metronome.IsPlaying && m.currentBeat == beatPosition && m.beatAnimation > 0 {
		// This gnome is lit up
		var color string
		if beatPosition == 1 {
			// First beat (downbeat) - special teal color
			color = "51" // Bright teal
		} else {
			// Other beats - bright yellow
			color = "226" // Bright yellow
		}
		
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(gnome)
	} else {
		// This gnome is dim - dark gray
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Dim gray
			Render(gnome)
	}
}

// getBeatGnomes returns gnomes for each beat of the time signature
func (m Model) getBeatGnomes() string {
	numBeats := m.metronome.TimeSignature.Beats
	gnomes := make([]string, numBeats)
	
	// Create a gnome for each beat position
	for i := 1; i <= numBeats; i++ {
		gnomes[i-1] = m.getGnome(i)
	}

	// Split all gnomes into lines
	gnomeLines := make([][]string, numBeats)
	for i, gnome := range gnomes {
		gnomeLines[i] = strings.Split(gnome, "\n")
	}

	// Combine lines horizontally with spacing
	result := ""
	maxLines := 4 // All gnomes have 4 lines
	
	for lineNum := 0; lineNum < maxLines; lineNum++ {
		line := ""
		
		for gnomeNum := 0; gnomeNum < numBeats; gnomeNum++ {
			// Add the gnome's line
			if lineNum < len(gnomeLines[gnomeNum]) {
				line += gnomeLines[gnomeNum][lineNum]
			} else {
				line += "     " // Empty space if gnome doesn't have this line
			}
			
			// Add spacing between gnomes (except after the last one)
			if gnomeNum < numBeats-1 {
				line += "  " // Two spaces between gnomes
			}
		}
		
		result += line
		if lineNum < maxLines-1 {
			result += "\n"
		}
	}

	return result
}

// initializeStars creates star positions throughout the terminal
func (m *Model) initializeStars() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	// Calculate number of stars based on terminal size - more stars!
	numStars := (m.width * m.height) / 8 // One star per 8 character cells (much denser!)
	if numStars > 300 {
		numStars = 300 // Cap at 300 stars for performance
	}
	if numStars < 50 {
		numStars = 50 // Minimum 50 stars for small terminals
	}

	m.starPositions = make([][]int, numStars)

	// Generate random positions for stars
	for i := 0; i < numStars; i++ {
		x := rand.Intn(m.width)
		y := rand.Intn(m.height)
		m.starPositions[i] = []int{x, y}
	}
}

// generateStarBackground creates a background with ASCII stars
func (m Model) generateStarBackground(content string) string {
	if len(m.starPositions) == 0 || m.width <= 0 || m.height <= 0 {
		return content
	}

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Ensure we have enough lines
	for len(lines) < m.height {
		lines = append(lines, "")
	}

	// Build the result with stars overlaid
	result := make([]string, len(lines))

	for y, line := range lines {
		resultLine := ""
		lineRunes := []rune(line)

		// Build the line character by character
		maxX := m.width
		if len(lineRunes) > maxX {
			maxX = len(lineRunes)
		}

		for x := 0; x < maxX; x++ {
			// Check if there's a star at this position
			starAtPosition := false
			var starChar string

			for i, pos := range m.starPositions {
				if pos[0] == x && pos[1] == y {
					// Check if this position is empty space
					if x >= len(lineRunes) || lineRunes[x] == ' ' {
						// Choose star color based on beat animation and star index
						colorIndex := i % len(m.starColors)

						// Make ALL stars flash together with the beat
						if m.beatAnimation > 0 {
							// ALL stars flash bright together on beat!
							if m.beatAnimation >= 4 {
								// Peak brightness - ALL stars use brightest color
								if m.currentBeat == 1 {
									// First beat: ALL stars flash the brightest white
									colorIndex = len(m.starColors) - 1
								} else {
									// Other beats: ALL stars flash bright yellow/green
									colorIndex = len(m.starColors) - 2
								}
							} else if m.beatAnimation >= 2 {
								// Medium fade - ALL stars dim together
								colorIndex = len(m.starColors) - 3
							} else {
								// Final fade - ALL stars dim further together
								colorIndex = len(m.starColors) - 4
							}
						} else if m.metronome.IsPlaying {
							// Between beats - ALL stars stay dim but visible
							colorIndex = 2 // Same dim color for all
						} else {
							// Stopped - ALL stars very dim
							colorIndex = 0 // Same very dim color for all
						}

						// Create twinkling star
						starColor := m.starColors[colorIndex]
						starChar = lipgloss.NewStyle().
							Foreground(lipgloss.Color(starColor)).
							Render("✦")
						starAtPosition = true
						break
					}
				}
			}

			if starAtPosition {
				resultLine += starChar
			} else if x < len(lineRunes) {
				resultLine += string(lineRunes[x])
			} else {
				resultLine += " "
			}
		}

		result[y] = resultLine
	}

	return strings.Join(result, "\n")
}

// renderMainWithBorder renders the main content with a solid border and star background
func (m Model) renderMainWithBorder() string {
	// Get the main content
	mainContent := m.renderMainContent()

	// Static yellow border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("226")).
		Padding(1, 2).
		Width(m.width - 10).
		Align(lipgloss.Center)

	// Wrap content in border
	borderedContent := borderStyle.Render(mainContent)

	// Center the bordered content
	centeredContent := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(borderedContent)

	// Add star background only (no garden decorations)
	return m.generateStarBackground(centeredContent)
}

// renderMainContent returns just the inner metronome content (without border)
func (m Model) renderMainContent() string {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		MarginBottom(1)

	bpmStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	beatStyle := lipgloss.NewStyle().
		Width(6).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center).
		MarginRight(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	// Title
	title := titleStyle.Render("🍄 Metrognome 🍄")

	// BPM display
	bpmDisplay := fmt.Sprintf("%d BPM", m.metronome.BPM)
	bpmLine := bpmStyle.Render(bpmDisplay)

	// Time signature
	tsDisplay := fmt.Sprintf("%s", m.metronome.TimeSignature.Name)

	// Beat visualization
	beats := ""
	for i := 1; i <= m.metronome.TimeSignature.Beats; i++ {
		style := beatStyle
		if i == m.currentBeat && m.metronome.IsPlaying {
			// Animate the current beat
			if m.beatAnimation > 0 {
				style = style.
					Background(lipgloss.Color("212")).
					Foreground(lipgloss.Color("231"))
			} else {
				style = style.
					Background(lipgloss.Color("240")).
					Foreground(lipgloss.Color("231"))
			}
		} else {
			style = style.
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("244"))
		}

		if i == 1 {
			beats += style.Render("𝟙")
		} else {
			beats += style.Render(fmt.Sprintf("%d", i))
		}
	}

	// Status
	status := "Press SPACE to start"
	if m.metronome.IsPlaying {
		status = "Playing... Press SPACE to stop"
	}
	statusLine := statusStyle.Render(status)

	// Sound status
	soundStatus := "🔇 Sound: OFF"
	if m.soundEnabled {
		soundStatus = "🔊 Sound: ON"
	}
	soundLine := statusStyle.Render(soundStatus)

	// Gnome saying
	saying := m.metronome.TimeSignature.GnomeSaying

	// BPM description
	bpmDesc := metronome.GetBPMDescription(m.metronome.BPM)

	// Beat counter gnomes
	gnomes := m.getBeatGnomes()
	
	// Pendulum swing arm
	swingArm := m.getSwingArm()

	// Compose the view
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		bpmLine,
		bpmDesc,
		"",
		tsDisplay,
		saying,
		"",
		beats,
		"",
		swingArm,
		"",
		gnomes,
		"",
		statusLine,
		soundLine,
		"",
	)

	// Help hint
	helpHint := m.help.View(m.keys)

	// Quit instruction
	quitInstruction := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Press 'q' or Ctrl+C to quit")

	// Add help and quit instruction at the bottom
	bottomContent := lipgloss.JoinVertical(
		lipgloss.Left,
		helpHint,
		quitInstruction,
	)

	// Final content
	return lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		"",
		bottomContent,
	)
}

// getSwingArm returns a visual pendulum that swings with the tempo
func (m Model) getSwingArm() string {
	// Create a simple pendulum visualization
	armLength := 15
	centerPos := armLength
	
	// Calculate pendulum position based on angle
	offset := int(math.Sin(m.pendulumAngle) * float64(armLength))
	pendulumPos := centerPos + offset
	
	// Build the pendulum line
	line := strings.Repeat(" ", armLength*2+1)
	lineRunes := []rune(line)
	
	// Add the pivot point
	if centerPos < len(lineRunes) {
		lineRunes[centerPos] = '┬'
	}
	
	// Add the pendulum bob
	if pendulumPos >= 0 && pendulumPos < len(lineRunes) {
		lineRunes[pendulumPos] = '●'
	}
	
	// Add connecting line if needed
	start := min(centerPos, pendulumPos)
	end := max(centerPos, pendulumPos)
	for i := start + 1; i < end; i++ {
		if i < len(lineRunes) && lineRunes[i] == ' ' {
			lineRunes[i] = '─'
		}
	}
	
	pendulumLine := string(lineRunes)
	
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Render(pendulumLine)
}

// addGardenDecorations adds a rich variety of garden decorations throughout the UI
func (m Model) addGardenDecorations(content string) string {
	if m.width <= 0 || m.height <= 0 {
		return content
	}
	
	lines := strings.Split(content, "\n")
	
	// Ensure we have enough lines
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	
	// Simpler decorations that work well with terminal rendering
	decorations := []string{"*", ".", "o", "+", "^", "~", "x", "#"}
	
	// Add decorations to corners and edges only (safer)
	decorationPositions := []struct{ row, col int }{
		// Corners (with better bounds checking)
		{1, 1}, {1, max(1, m.width-2)}, 
		{max(1, m.height-2), 1}, {max(1, m.height-2), max(1, m.width-2)},
		// A few edge positions
		{0, m.width/4}, {0, 3*m.width/4},
		{max(1, m.height-1), m.width/4}, {max(1, m.height-1), 3*m.width/4},
	}
	
	// Place decorations safely
	for i, pos := range decorationPositions {
		row, col := pos.row, pos.col
		
		// Strict bounds checking
		if row >= 0 && row < len(lines) && col >= 0 && col < m.width-1 {
			// Ensure line is long enough
			for len(lines[row]) <= col {
				lines[row] += " "
			}
			
			// Convert to runes for safe manipulation
			lineRunes := []rune(lines[row])
			if col < len(lineRunes) && lineRunes[col] == ' ' {
				decoration := decorations[i%len(decorations)]
				lineRunes[col] = rune(decoration[0])
				lines[row] = string(lineRunes)
			}
		}
	}
	
	return strings.Join(lines, "\n")
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// playSound plays a system sound based on the OS
func playSound(isFirstBeat bool) {
	switch runtime.GOOS {
	case "darwin": // macOS
		if isFirstBeat {
			// Play the regular beat sound + an extra blip for first beat
			go exec.Command("afplay", "/System/Library/Sounds/Tink.aiff").Run() // Regular sound
			go exec.Command("afplay", "/System/Library/Sounds/Pop.aiff").Run()  // Extra blip
		} else {
			exec.Command("afplay", "/System/Library/Sounds/Tink.aiff").Run()
		}
	case "linux":
		if isFirstBeat {
			// Play regular beat + extra higher blip
			go func() {
				exec.Command("beep", "-f", "440", "-l", "50").Run() // Regular sound
			}()
			go func() {
				time.Sleep(10 * time.Millisecond) // Slight delay
				exec.Command("beep", "-f", "880", "-l", "30").Run() // Extra blip
			}()
		} else {
			// Regular beat
			if err := exec.Command("beep", "-f", "440", "-l", "50").Run(); err != nil {
				exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/message.oga").Run()
			}
		}
	case "windows":
		if isFirstBeat {
			// Play regular beat + extra blip
			go func() {
				exec.Command("powershell", "-c", "[console]::beep(800,100)").Run() // Regular sound
			}()
			go func() {
				time.Sleep(50 * time.Millisecond) // Slight delay
				exec.Command("powershell", "-c", "[console]::beep(1200,50)").Run() // Extra blip
			}()
		} else {
			exec.Command("powershell", "-c", "[console]::beep(800,100)").Run()
		}
	default:
		// Fallback to terminal bell
		if isFirstBeat {
			fmt.Print("\a\a") // Double bell for first beat
		} else {
			fmt.Print("\a")
		}
	}
}
