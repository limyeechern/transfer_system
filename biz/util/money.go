package util

import (
	"errors"
	"strconv"
	"strings"
)

const AmountScale = int64(100000)

var ErrInvalidAmount = errors.New("invalid amount")

func ParseAmount5DP(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, ErrInvalidAmount
	}

	parts := strings.Split(amount, ".")
	if len(parts) > 2 {
		return 0, ErrInvalidAmount
	}

	whole := parts[0]
	if whole == "" {
		whole = "0"
	}
	if !allDigits(whole) {
		return 0, ErrInvalidAmount
	}

	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
		if len(fraction) > 5 || !allDigits(fraction) {
			return 0, ErrInvalidAmount
		}
	}

	fraction += strings.Repeat("0", 5-len(fraction))
	scaled, err := strconv.ParseInt(whole+fraction, 10, 64)
	if err != nil {
		return 0, ErrInvalidAmount
	}
	return scaled, nil
}

func allDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
