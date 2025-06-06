# 🍄 Metrognome 🍄

A delightful terminal-based metronome with a garden gnome theme! Keep perfect time while enjoying whimsical gnome wisdom and garden-fresh beats.

## Features

- 🎵 Variable BPM (20-300) with gnome-themed tempo descriptions
- 🎼 Multiple time signatures (4/4, 3/4, 6/8, 5/4, 7/8, 2/4)
- 🌻 Garden-themed presets for common rhythms
- 🎨 Beautiful terminal UI powered by Bubble Tea
- 🧙 Animated gnome companion that dances to the beat
- 🌱 Each time signature comes with its own gnome saying

## Installation

### Method 1: Direct Install (Recommended)

```bash
go install github.com/drj613/metrognome@latest
```

### Method 2: Clone and Build

If the direct install fails, clone and build locally:

```bash
git clone https://github.com/drj613/metrognome.git
cd metrognome
go mod tidy
go build -o metrognome
```

### Troubleshooting

If you encounter module errors with `go install`, try:

1. **Clear module cache:**
   ```bash
   go clean -modcache
   go install github.com/drj613/metrognome@latest
   ```

2. **Use specific version:**
   ```bash
   go install github.com/drj613/metrognome@main
   ```

3. **Clone method (always works):**
   ```bash
   git clone https://github.com/drj613/metrognome.git
   cd metrognome
   go run .
   ```

## Usage

Simply run:

```bash
./metrognome
```

### Controls

- **Space**: Start/Stop the metronome
- **↑/↓** or **k/j**: Increase/Decrease BPM by 5
- **Tab**: Cycle through time signatures
- **p**: Show preset rhythms
- **?**: Show help
- **q**: Quit

### Presets

Choose from gnome-approved presets:
- 🚶 Peaceful Garden Stroll (60 BPM, 4/4)
- 🎵 Gnome Work Song (120 BPM, 4/4)
- 💃 Toadstool Waltz (90 BPM, 3/4)
- 🏃 Pixie Dust Presto (180 BPM, 4/4)
- 🕺 Underground Jig (140 BPM, 6/8)
- 🧘 Meditation by the Pond (40 BPM, 4/4)

## Building from Source

Requirements:
- Go 1.21 or later

```bash
go mod download
go build -o metrognome
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## License

MIT

---

*"A gnome without rhythm is like a garden without flowers!"* 🌻
