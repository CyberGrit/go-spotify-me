package theme

import (
	"strings"
	"testing"
)

func TestTruncateOrPad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{"Empty string, width 5", "", 5, "     "},
		{"Width 0", "hello", 0, ""},
		{"Negative width", "hello", -1, ""},
		{"Exact width", "hello", 5, "hello"},
		{"Shorter string", "hi", 5, "hi   "},
		{"Longer string, width > 3", "hello world", 8, "hello..."},
		{"Longer string, width <= 3", "hi", 2, "hi"},
		{"Longer string, width <= 3 edge", "hey", 2, "he"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TruncateOrPad(tc.input, tc.width)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestRenderRow(t *testing.T) {
	cells := []string{"1", "Test Song", "Artist Name"}
	widths := []int{2, 10, 15}

	result := RenderRow(cells, widths, RowStyle)

	// We check for the presence of the values and separator since lipgloss
	// adds ANSI escape codes for styling which are hard to string match exactly.
	if !strings.Contains(result, "1 ") {
		t.Errorf("Result should contain padded '1 ': %s", result)
	}
	if !strings.Contains(result, "Test Song ") {
		t.Errorf("Result should contain padded 'Test Song ': %s", result)
	}
	if !strings.Contains(result, "Artist Name    ") {
		t.Errorf("Result should contain padded 'Artist Name    ': %s", result)
	}
	if !strings.Contains(result, " | ") {
		t.Errorf("Result should contain column separators: %s", result)
	}
}
