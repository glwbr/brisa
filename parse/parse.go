// Package parse provides common parsing utilities for Brazilian invoice data.
package parse

import (
	"strconv"
	"strings"
	"time"
)

// Digits extracts only numeric characters from a string.
func Digits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Text normalizes whitespace: replaces NBSP with space, trims, and collapses multiple spaces.
func Text(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.Join(strings.Fields(s), " ")
}

// Quantity parses a Brazilian-formatted number (1.234,56) to float64.
func Quantity(s string) float64 {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Int parses a string to int, returning 0 on failure.
func Int(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// Percent parses a percentage string (e.g., "18,5%" or "18.5") to float64.
func Percent(s string) float64 {
	s = strings.ReplaceAll(s, "%", "")
	return Quantity(s)
}

// BrazilianDate parses common Brazilian date formats.
func BrazilianDate(s string) (time.Time, error) {
	s = Text(s)
	layouts := []string{
		"02/01/2006 15:04:05-07:00",
		"02/01/2006 15:04:05",
		"02/01/2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Layout: "BR date", Value: s}
}

// FirstNonEmpty returns the first non-empty string from the arguments.
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
