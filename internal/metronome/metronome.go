package metronome

import (
	"time"
)

// TimeSignature represents a musical time signature
type TimeSignature struct {
	Beats       int    // Number of beats per measure
	BeatValue   int    // Note value that gets the beat (4 = quarter note, 8 = eighth note)
	Name        string // Human-readable name
	GnomeSaying string // Fun gnome-themed description
}

// Preset represents a metronome preset configuration
type Preset struct {
	Name          string
	BPM           int
	TimeSignature TimeSignature
	Description   string
}

// Metronome represents the core metronome logic
type Metronome struct {
	BPM           int
	TimeSignature TimeSignature
	IsPlaying     bool
	CurrentBeat   int
	ticker        *time.Ticker
	beatChan      chan int
}

// CommonTimeSignatures provides preset time signatures with gnome themes
var CommonTimeSignatures = []TimeSignature{
	{
		Beats:       4,
		BeatValue:   4,
		Name:        "4/4 - Garden March",
		GnomeSaying: "Four steady steps through the garden path!",
	},
	{
		Beats:       3,
		BeatValue:   4,
		Name:        "3/4 - Gnome Waltz",
		GnomeSaying: "Dance among the toadstools, one-two-three!",
	},
	{
		Beats:       6,
		BeatValue:   8,
		Name:        "6/8 - Fairy Ring Jig",
		GnomeSaying: "Six quick hops around the mushroom circle!",
	},
	{
		Beats:       5,
		BeatValue:   4,
		Name:        "5/4 - Mystical Garden",
		GnomeSaying: "Five beats for the ancient gnome rituals!",
	},
	{
		Beats:       7,
		BeatValue:   8,
		Name:        "7/8 - Gnome's Riddle",
		GnomeSaying: "Seven steps to solve the garden mystery!",
	},
	{
		Beats:       2,
		BeatValue:   4,
		Name:        "2/4 - Quick March",
		GnomeSaying: "Left-right through the gnome village!",
	},
}

// CommonPresets provides common BPM and time signature combinations
var CommonPresets = []Preset{
	{
		Name:          "Peaceful Garden Stroll",
		BPM:           60,
		TimeSignature: CommonTimeSignatures[0], // 4/4
		Description:   "A leisurely walk through the gnome gardens",
	},
	{
		Name:          "Gnome Work Song",
		BPM:           120,
		TimeSignature: CommonTimeSignatures[0], // 4/4
		Description:   "Perfect for tending to the mushroom patches",
	},
	{
		Name:          "Toadstool Waltz",
		BPM:           90,
		TimeSignature: CommonTimeSignatures[1], // 3/4
		Description:   "Dance beneath the moonlit mushrooms",
	},
	{
		Name:          "Pixie Dust Presto",
		BPM:           180,
		TimeSignature: CommonTimeSignatures[0], // 4/4
		Description:   "When the garden gnomes need to hurry!",
	},
	{
		Name:          "Underground Jig",
		BPM:           140,
		TimeSignature: CommonTimeSignatures[2], // 6/8
		Description:   "For celebrating in the gnome tunnels",
	},
	{
		Name:          "Meditation by the Pond",
		BPM:           40,
		TimeSignature: CommonTimeSignatures[0], // 4/4
		Description:   "Slow and steady wins the gnome race",
	},
}

// New creates a new Metronome instance
func New(bpm int, timeSignature TimeSignature) *Metronome {
	return &Metronome{
		BPM:           bpm,
		TimeSignature: timeSignature,
		IsPlaying:     false,
		CurrentBeat:   1,
		beatChan:      make(chan int, 1),
	}
}

// Start begins the metronome
func (m *Metronome) Start() {
	if m.IsPlaying {
		return
	}

	m.IsPlaying = true
	m.CurrentBeat = 1

	// Calculate interval between beats
	interval := time.Duration(60000/m.BPM) * time.Millisecond
	m.ticker = time.NewTicker(interval)

	go func() {
		for range m.ticker.C {
			if m.IsPlaying {
				select {
				case m.beatChan <- m.CurrentBeat:
				default:
					// Channel is full, skip this beat
				}
				
				m.CurrentBeat++
				if m.CurrentBeat > m.TimeSignature.Beats {
					m.CurrentBeat = 1
				}
			}
		}
	}()
}

// Stop halts the metronome
func (m *Metronome) Stop() {
	if !m.IsPlaying {
		return
	}

	m.IsPlaying = false
	if m.ticker != nil {
		m.ticker.Stop()
	}
	m.CurrentBeat = 1
}

// SetBPM changes the tempo
func (m *Metronome) SetBPM(bpm int) {
	if bpm < 20 || bpm > 300 {
		return
	}

	wasPlaying := m.IsPlaying
	if wasPlaying {
		m.Stop()
	}

	m.BPM = bpm

	if wasPlaying {
		m.Start()
	}
}

// SetTimeSignature changes the time signature
func (m *Metronome) SetTimeSignature(ts TimeSignature) {
	wasPlaying := m.IsPlaying
	if wasPlaying {
		m.Stop()
	}

	m.TimeSignature = ts
	m.CurrentBeat = 1

	if wasPlaying {
		m.Start()
	}
}

// BeatChannel returns the channel that emits beat numbers
func (m *Metronome) BeatChannel() <-chan int {
	return m.beatChan
}

// GetBPMDescription returns a gnome-themed description of the current tempo
func GetBPMDescription(bpm int) string {
	switch {
	case bpm < 40:
		return "Gnome hibernation speed"
	case bpm < 60:
		return "Sleepy garden gnome pace"
	case bpm < 80:
		return "Morning dew collection tempo"
	case bpm < 100:
		return "Casual mushroom picking rhythm"
	case bpm < 120:
		return "Standard gnome work tempo"
	case bpm < 140:
		return "Energetic garden tending speed"
	case bpm < 160:
		return "Gnome celebration dance"
	case bpm < 180:
		return "Chasing garden pests tempo"
	case bpm < 200:
		return "Gnome emergency response speed"
	default:
		return "Hyperactive pixie dust overdose!"
	}
}