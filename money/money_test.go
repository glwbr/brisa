package money

import "testing"

func TestParseBRL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  BRL
	}{
		{name: "plain", input: "527,84", want: FromFloat(527.84)},
		{name: "with_symbol_and_thousands", input: "R$ 1.234,56", want: FromFloat(1234.56)},
		{name: "negative", input: "-10,00", want: FromFloat(-10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("Parse(%q) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestBRLString(t *testing.T) {
	got := FromFloat(1234.56).String()
	const want = "R$ 1.234,56"
	if got != want {
		t.Fatalf("String() = %q; want %q", got, want)
	}
}
