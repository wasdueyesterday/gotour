package gotour

import (
	"testing"
)

func TestLongestSubStr(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected int
	} {
		// {"empty string", "", 0},
		// {"all same", "bbbbb", 1}, 
		{"standard", "abcabcbb", 3},
		// {"standard2", "pwwkew", 3},
		// {"single char", "a", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T)  {
			actual := LongestNorepeatSubStr(tt.input)
			if actual != tt.expected {
				t.Errorf("input %q: expected %d, got %d", tt.input, tt.expected, actual)
			}
		})
	}
}