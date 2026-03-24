package money

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func ParseToMinorUnit(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, fmt.Errorf("amount cannot be empty")
	}
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return 0, fmt.Errorf("failed to parse amount '%s': %w", amount, err)
	}
	minorUnit := d.Shift(2).Round(0)
	return minorUnit.IntPart(), nil
}

func FormatFromMinorUnit(amount int64) string {
	d := decimal.NewFromInt(amount).Shift(-2).StringFixed(2)
	return d
}
