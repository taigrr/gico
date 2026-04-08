package common

import (
	"testing"

	sc "github.com/taigrr/simplecolorpalettes/simplecolor"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		wantMin int
		wantMax int
	}{
		{
			name:    "normal values",
			input:   []int{3, 1, 4, 1, 5, 9, 2, 6},
			wantMin: 1,
			wantMax: 9,
		},
		{
			name:    "with zeros",
			input:   []int{0, 3, 0, 7, 0, 2},
			wantMin: 2,
			wantMax: 7,
		},
		{
			name:    "all zeros",
			input:   []int{0, 0, 0},
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "single element",
			input:   []int{5},
			wantMin: 5,
			wantMax: 5,
		},
		{
			name:    "single zero",
			input:   []int{0},
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "all same non-zero",
			input:   []int{4, 4, 4},
			wantMin: 4,
			wantMax: 4,
		},
		{
			name:    "empty slice",
			input:   []int{},
			wantMin: 0,
			wantMax: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := MinMax(tt.input)
			if min != tt.wantMin {
				t.Errorf("min: expected %d, got %d", tt.wantMin, min)
			}
			if max != tt.wantMax {
				t.Errorf("max: expected %d, got %d", tt.wantMax, max)
			}
		})
	}
}

func TestColorForFrequency(t *testing.T) {
	// freq=0 should always return SimpleColor(0)
	c := ColorForFrequency(0, 1, 10)
	if c != sc.SimpleColor(0) {
		t.Errorf("expected SimpleColor(0) for freq=0, got %v", c)
	}

	// Non-zero freq with small spread should return a valid color
	c = ColorForFrequency(1, 1, 3)
	if c == sc.SimpleColor(0) {
		t.Error("expected non-zero color for freq=1 with small spread")
	}

	// Non-zero freq with large spread
	c = ColorForFrequency(50, 1, 100)
	if c == sc.SimpleColor(0) {
		t.Error("expected non-zero color for freq=50 in range 1-100")
	}

	// Max frequency should return the last color
	c = ColorForFrequency(100, 1, 100)
	if c == sc.SimpleColor(0) {
		t.Error("expected non-zero color for max frequency")
	}

	// Min frequency should return first real color
	c = ColorForFrequency(1, 1, 100)
	if c == sc.SimpleColor(0) {
		t.Error("expected non-zero color for min frequency")
	}
}

func TestCreateGraph(t *testing.T) {
	buf := CreateGraph()
	if buf.Len() != 0 {
		t.Errorf("expected empty buffer, got %d bytes", buf.Len())
	}
}

func TestSetColorScheme(t *testing.T) {
	// Save original length
	origLen := len(colorScheme)

	// Reset after test
	defer func() {
		colorScheme = colorScheme[:origLen]
	}()

	// SetColorScheme should append colors
	SetColorScheme(nil)
	if len(colorScheme) != origLen {
		t.Errorf("expected length %d after nil input, got %d", origLen, len(colorScheme))
	}
}
