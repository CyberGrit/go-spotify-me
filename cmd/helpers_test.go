package cmd

import (
	"reflect"
	"testing"
)

func TestCalculateColumnWidths(t *testing.T) {
	tests := []struct {
		name       string
		totalWidth int
		ratios     []float64
		expected   []int
	}{
		{
			name:       "Equal ratios",
			totalWidth: 100,
			ratios:     []float64{0.5, 0.5},
			expected:   []int{48, 49}, // 100 - 3 (padding for 2 cols) = 97. 97 * 0.5 = 48. Remainder goes to last col: 97 - 48 = 49.
		},
		{
			name:       "Three columns",
			totalWidth: 100,
			ratios:     []float64{0.2, 0.3, 0.5},
			expected:   []int{18, 28, 48}, // padding = 3 * 2 = 6. 94 available. 18, 28, 48
		},
		{
			name:       "One column",
			totalWidth: 50,
			ratios:     []float64{1.0},
			expected:   []int{50},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateColumnWidths(tc.totalWidth, tc.ratios)
			
			// Verify exact expected widths
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
			
			// Verify total sum equals totalWidth - padding
			padding := 3 * (len(tc.ratios) - 1)
			expectedSum := tc.totalWidth - padding
			sum := 0
			for _, w := range result {
				sum += w
			}
			
			if sum != expectedSum {
				t.Errorf("Expected sum %d, got %d", expectedSum, sum)
			}
		})
	}
}
