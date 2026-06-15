package util

import (
	"fmt"
	"strconv"
	"strings"

	"transfer_system/biz/apperror"
)

const AmountScale = int64(100000)

const AmountDecimalPlaces = 5

func ParseAmount5DP(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, apperror.ErrInvalidAmount
	}

	parts := strings.Split(amount, ".")
	if len(parts) > 2 {
		return 0, apperror.ErrInvalidAmount
	}

	whole := parts[0]
	if whole == "" {
		whole = "0"
	}
	if !allDigits(whole) {
		return 0, apperror.ErrInvalidAmount
	}

	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
		if len(fraction) > 5 || !allDigits(fraction) {
			return 0, apperror.ErrInvalidAmount
		}
	}

	fraction += strings.Repeat("0", 5-len(fraction))
	scaled, err := strconv.ParseInt(whole+fraction, 10, 64)
	if err != nil {
		return 0, apperror.ErrInvalidAmount
	}
	return scaled, nil
}

func FormatAmount5DP(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	whole := amount / AmountScale
	fraction := amount % AmountScale
	return fmt.Sprintf("%s%d.%0*d", sign, whole, AmountDecimalPlaces, fraction)
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
