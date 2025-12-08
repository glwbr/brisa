package money

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// BRL represents a Brazilian Real value with cents as the base unit
// For example, R$ 1.23 is represented as 123 cents internally
type BRL int64

var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrInvalidFormat  = errors.New("invalid BRL format")
)

// FromFloat converts a float64 to BRL, rounding it to maintain precision
func FromFloat(f float64) BRL {
	return BRL(math.Round(f * 100))
}

// Float64 converts a BRL value to float64 (in reais)
func (b BRL) Float64() float64 {
	return float64(b) / 100
}

// Add returns the sum of two BRL values
func (b BRL) Add(other BRL) BRL {
	return b + other
}

// Sub returns the difference between two BRL values
func (b BRL) Sub(other BRL) BRL {
	return b - other
}

// Mul multiplies a BRL value by a factor
func (b BRL) Mul(factor float64) BRL {
	return BRL(math.Round(float64(b) * factor))
}

// Div divides a BRL value by a divisor
func (b BRL) Div(divisor float64) (BRL, error) {
	if divisor == 0 {
		return 0, ErrDivisionByZero
	}
	return BRL(math.Round(float64(b) / divisor)), nil
}

// String returns the formatted BRL value (e.g., "R$ 1.234,56")
func (b BRL) String() string {
	value := abs(int64(b))
	reais := value / 100
	cents := value % 100

	sign := ""
	if b < 0 {
		sign = "-"
	}

	reaisStr := formatThousands(reais)
	return fmt.Sprintf("R$ %s%s,%02d", sign, reaisStr, cents)
}

// Parse converts a string representation of BRL to a BRL value
func Parse(s string) (BRL, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "R$", "")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.Replace(s, ",", ".", 1)

	if strings.Count(s, ".") > 1 {
		return 0, fmt.Errorf("%w: too many decimal points", ErrInvalidFormat)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidFormat, err)
	}

	return FromFloat(f), nil
}

// Ratio calculates the ratio of this BRL value to another
func (b BRL) Ratio(other BRL) (float64, error) {
	if other == 0 {
		return 0, ErrDivisionByZero
	}
	return float64(b) / float64(other), nil
}

// Abs returns the absolute value of the BRL amount
func (b BRL) Abs() BRL {
	if b < 0 {
		return -b
	}
	return b
}

// abs returns the absolute value of x
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// formatThousands formats an integer with thousand separators
func formatThousands(n int64) string {
	s := fmt.Sprintf("%d", n)
	var parts []string

	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}

	if s != "" {
		parts = append([]string{s}, parts...)
	}

	return strings.Join(parts, ".")
}
