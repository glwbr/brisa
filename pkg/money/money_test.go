package money

import (
	"math"
	"testing"
)

func TestFromFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected BRL
	}{
		{"Positive value", 1.23, 123},
		{"Small positive value", 0.01, 1},
		{"Float with many decimals", 10.345, 1035},
		{"Large value", 10000.00, 1000000},
		{"Negative value", -10.99, -1099},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromFloat(tt.input)
			if got != tt.expected {
				t.Errorf("FromFloat(%f) = %d; want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected float64
	}{
		{"Positive value", FromFloat(1.23), 1.23},
		{"Small positive value", FromFloat(0.01), 0.01},
		{"Decimal value", FromFloat(10.345), 10.35},
		{"Negative value", FromFloat(-10.99), -10.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Float64()
			if got != tt.expected {
				t.Errorf("Float64() = %f; want %f", got, tt.expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected BRL
	}{
		{"Add", FromFloat(10.50), FromFloat(5.25), FromFloat(15.75)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Add(tt.b)
			if got != tt.expected {
				t.Errorf("Add: got %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected BRL
	}{
		{"Subtract", FromFloat(10.50), FromFloat(5.25), FromFloat(5.25)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Sub(tt.b)
			if got != tt.expected {
				t.Errorf("Sub: got %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		a        BRL
		mult     float64
		expected BRL
	}{
		{"Multiply", FromFloat(10.50), 2, FromFloat(21.00)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mul(tt.mult)
			if got != tt.expected {
				t.Errorf("Mul: got %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		name     string
		a        BRL
		div      float64
		expected BRL
	}{
		{"Divide", FromFloat(10.50), 2, FromFloat(5.25)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.a.Div(tt.div)
			if got != tt.expected {
				t.Errorf("Div: got %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		name     string
		dividend BRL
		divisor  float64
		wantErr  bool
	}{
		{"Division by zero", BRL(100), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dividend.Div(tt.divisor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div(%v, %v) error = %v, wantErr %v", tt.dividend, tt.divisor, err, tt.wantErr)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected string
	}{
		{"Simple", FromFloat(10.99), "R$ 10,99"},
		{"Thousands", FromFloat(1000.50), "R$ 1.000,50"},
		{"Millions", FromFloat(1000000), "R$ 1.000.000,00"},
		{"Negative", FromFloat(-42.75), "R$ -42,75"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.expected {
				t.Errorf("String: got %s; want %s", got, tt.expected)
			}
		})
	}
}

func TestParseSuccess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BRL
	}{
		{"Basic", "10,99", FromFloat(10.99)},
		{"With prefix", "R$ 10,99", FromFloat(10.99)},
		{"Thousands", "1.000,50", FromFloat(1000.50)},
		{"Millions", "R$ 1.234.567,89", FromFloat(1234567.89)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("Parse(%q) = %d; want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseFailure(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"Invalid grouping", "R$ 1,000,00", true},
		{"No decimals", "R$ 1000", false},
		{"Extra spacing and malformed", "  R$ 1.000,0 0 ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
			}
		})
	}
}

func TestMultiplicationPrecision(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		mult     float64
		expected BRL
	}{
		{"0.29 × (10 / 3) → 0.97", FromFloat(0.29), 10.0 / 3.0, FromFloat(0.97)},
		{"10.01 × 1.1 → 11.01", FromFloat(10.01), 1.1, FromFloat(11.01)},
		{"0.07 × (10/7) × 0.7 → 0.07", BRL(7), (10.0 / 7.0) * 0.7, BRL(7)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Mul(tt.mult)
			if result != tt.expected {
				t.Errorf("Mul: got %v (%d cents); want %v (%d cents)",
					result, int64(result), tt.expected, int64(tt.expected))
			}
		})
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected float64
		wantErr  bool
	}{
		{"10 / 2 = 5", FromFloat(10.00), FromFloat(2.00), 5.0, false},
		{"9 / 3 = 3", FromFloat(9.00), FromFloat(3.00), 3.0, false},
		{"same values = 1", FromFloat(5.25), FromFloat(5.25), 1.0, false},
		{"smaller numerator = < 1", FromFloat(1.00), FromFloat(2.00), 0.5, false},
		{"division by zero", FromFloat(1.00), BRL(0), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Ratio(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ratio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.expected) > 0.0001 {
				t.Errorf("Ratio() = %.4f; want %.4f", got, tt.expected)
			}
		})
	}
}
