package main

import (
	"testing"
)

type testCase struct {
	Input  string
	Output string
}

func TestNormalizePhoneNumber(t *testing.T) {
	numberData := []testCase{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"123 456 7892", "1234567892"},
		{"123 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"123-456-7892", "1234567892"},
	}
	for _, num := range numberData {
		t.Run(num.Input, func(t *testing.T) {
			actual := normalizePhoneNumber(num.Input)
			if actual != num.Output {
				t.Errorf("value pass as parameter %s and expected value %s", actual, num.Output)
			}
		})
	}
}
