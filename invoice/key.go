package invoice

// IsValidAccessKey checks if the access key is valid.
// It checks length, numeric content, and the check digit (DV).
func IsValidAccessKey(key string) bool {
	if len(key) != 44 {
		return false
	}

	for _, r := range key {
		if r < '0' || r > '9' {
			return false
		}
	}

	allSame := true
	for i := 1; i < len(key); i++ {
		if key[i] != key[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	total := 0
	multiplier := 2

	for i := 42; i >= 0; i-- {
		digit := int(key[i] - '0')
		total += digit * multiplier
		multiplier++
		if multiplier > 9 {
			multiplier = 2
		}
	}

	remainder := total % 11
	dv := 11 - remainder
	if remainder == 0 || remainder == 1 {
		dv = 0
	}

	expectedDV := int(key[43] - '0')
	return dv == expectedDV
}
