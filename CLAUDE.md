# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Metrognome is a terminal-based CLI metronome written in Go, featuring a whimsical garden gnome theme. It uses the Bubble Tea framework for creating an interactive terminal UI with support for variable BPM, multiple time signatures, and preset configurations.

## Key Commands

### Development
```bash
# Install dependencies
go mod tidy

# Run the application
go run .

# Build the binary
go build -o metrognome

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Run with race detector during development
go run -race .
```

### Testing Individual Packages
```bash
# Test metronome logic
go test ./internal/metronome

# Test UI components
go test ./internal/ui

# Test with verbose output
go test -v ./...
```

## Architecture

### Package Structure
- `main.go`: Entry point, delegates to cmd package
- `cmd/metrognome/`: Command-line interface setup using Bubble Tea
- `internal/metronome/`: Core metronome logic, timing, and musical concepts
  - Manages time signatures, BPM, presets
  - Handles beat timing via Go channels
  - Contains all musical domain logic
- `internal/ui/`: Bubble Tea UI implementation
  - `model.go`: Main UI model implementing tea.Model interface
  - Handles user input, screen rendering, and state management
  - Manages view states (main, preset selection, help)

### Key Design Patterns
1. **Channel-based timing**: The metronome uses Go channels for precise beat timing, avoiding timer drift
2. **State machine UI**: The UI uses Bubble Tea's Elm-inspired architecture with distinct states
3. **Separation of concerns**: Musical logic is isolated from UI rendering
4. **Theme consistency**: Garden gnome theme is maintained throughout via constants and helper functions

### Important Interfaces
- `metronome.Metronome`: Core metronome functionality with Start(), Stop(), SetBPM(), SetTimeSignature()
- `tea.Model`: Bubble Tea model interface implemented by the UI

### Styling Approach
- Uses Lip Gloss for consistent terminal styling
- Color scheme based on garden/nature theme (greens, browns, earth tones)
- ASCII art animations for visual feedback

## Project Conventions

### Gnome Theme Guidelines
- All user-facing strings should maintain the garden gnome theme
- Use puns and whimsical language where appropriate
- Beat descriptions should reference garden/gnome activities
- Error messages should be helpful but maintain character

### Code Style
- Follow standard Go conventions
- Keep functions focused and testable
- Use meaningful variable names that reflect the gnome theme where it enhances readability
- Comment complex timing logic