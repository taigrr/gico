package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestShiftSelectionByWeeks(t *testing.T) {
	tests := []struct {
		name         string
		year         int
		selected     int
		weeks        int
		wantYear     int
		wantSelected int
	}{
		{
			name:         "move left within year",
			year:         2025,
			selected:     21,
			weeks:        -1,
			wantYear:     2025,
			wantSelected: 14,
		},
		{
			name:         "move left across year boundary",
			year:         2025,
			selected:     3,
			weeks:        -1,
			wantYear:     2024,
			wantSelected: 362,
		},
		{
			name:         "move right across year boundary",
			year:         2025,
			selected:     361,
			weeks:        1,
			wantYear:     2026,
			wantSelected: 3,
		},
		{
			name:         "move right from leap day week into next week",
			year:         2024,
			selected:     58,
			weeks:        1,
			wantYear:     2024,
			wantSelected: 65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYear, gotSelected := shiftSelectionByWeeks(tt.year, tt.selected, tt.weeks)
			if gotYear != tt.wantYear || gotSelected != tt.wantSelected {
				t.Fatalf("shiftSelectionByWeeks(%d, %d, %d) = (%d, %d), want (%d, %d)", tt.year, tt.selected, tt.weeks, gotYear, gotSelected, tt.wantYear, tt.wantSelected)
			}
		})
	}
}

func TestGraphUpdateMovesAcrossYears(t *testing.T) {
	graph := Graph{Year: 2025, Selected: 3}

	updatedModel, _ := graph.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	updated := updatedModel

	if updated.Year != 2024 {
		t.Fatalf("expected year 2024 after moving left, got %d", updated.Year)
	}
	if updated.Selected != 362 {
		t.Fatalf("expected selection 362 after moving left, got %d", updated.Selected)
	}
}

func TestGraphUpdateStopsAtColumnBounds(t *testing.T) {
	graph := Graph{Year: 2025, Selected: 0}

	updatedModel, _ := graph.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	updated := updatedModel
	if updated.Selected != 0 {
		t.Fatalf("expected selection to stay at 0 when moving up, got %d", updated.Selected)
	}

	updatedModel, _ = graph.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	updated = updatedModel
	if updated.Selected != 1 {
		t.Fatalf("expected selection 1 after moving down, got %d", updated.Selected)
	}
}
