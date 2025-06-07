package types

import (
	"github.com/shopspring/decimal"
	"math"
)

type Symbol string

const (
	DOT  Symbol = "DOT"
	KSM  Symbol = "KSM"
	WND  Symbol = "WND"
	UNIT Symbol = "UNIT"
)

func (s Symbol) String() string {
	return string(s)
}

func IsValidSymbol(str string) bool {
	return str == string(DOT) || str == string(KSM) || str == string(WND) || str == string(UNIT)
}

// HumanReadableToNative converts asset value from human-readable to the native type
func HumanReadableToNative(hr decimal.Decimal, symbol Symbol) int64 {
	switch symbol {
	case DOT:
		return hr.Mul(decimal.NewFromFloat(math.Pow10(10))).IntPart()
	case KSM:
		return hr.Mul(decimal.NewFromFloat(math.Pow10(12))).IntPart()
	case WND:
		return hr.Mul(decimal.NewFromFloat(math.Pow10(12))).IntPart()
	case UNIT:
		return hr.Mul(decimal.NewFromFloat(math.Pow10(12))).IntPart()
	}
	return -1
}

// NativeToHumanReadable converts asset value from native to a human-readable format
func NativeToHumanReadable(native decimal.Decimal, symbol Symbol) decimal.Decimal {
	switch symbol {
	case DOT:
		return native.Div(decimal.NewFromFloat(math.Pow10(10)))
	case KSM:
		return native.Div(decimal.NewFromFloat(math.Pow10(12)))
	case WND:
		return native.Div(decimal.NewFromFloat(math.Pow10(12)))
	case UNIT:
		return native.Div(decimal.NewFromFloat(math.Pow10(12)))
	}
	return decimal.Zero
}
