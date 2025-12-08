package parse

import "testing"

func TestDigits(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"123.456.789-00", "12345678900"},
		{"  12 34  ", "1234"},
		{"abc123def", "123"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Digits(tt.input); got != tt.want {
			t.Errorf("Digits(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestQuantity(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"1.234,56", 1234.56},
		{"1,5", 1.5},
		{"100", 100},
		{"", 0},
	}
	for _, tt := range tests {
		if got := Quantity(tt.input); got != tt.want {
			t.Errorf("Quantity(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  hello   world  ", "hello world"},
		{"hello\u00A0world", "hello world"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Text(tt.input); got != tt.want {
			t.Errorf("Text(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := FirstNonEmpty("", "  ", "hello", "world"); got != "hello" {
		t.Errorf("FirstNonEmpty() = %q, want %q", got, "hello")
	}
	if got := FirstNonEmpty("", "", ""); got != "" {
		t.Errorf("FirstNonEmpty() = %q, want empty", got)
	}
}
