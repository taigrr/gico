package main

import "testing"

func TestParseYear(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fallback int
		want     int
	}{
		{name: "valid year", input: "2024", fallback: 2026, want: 2024},
		{name: "empty uses fallback", input: "", fallback: 2026, want: 2026},
		{name: "invalid uses fallback", input: "nope", fallback: 2026, want: 2026},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseYear(tt.input, tt.fallback)
			if got != tt.want {
				t.Fatalf("parseYear(%q, %d) = %d, want %d", tt.input, tt.fallback, got, tt.want)
			}
		})
	}
}
