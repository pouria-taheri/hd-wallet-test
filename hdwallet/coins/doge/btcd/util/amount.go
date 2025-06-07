package util

import (
	"github.com/shopspring/decimal"
	"math/big"
)

type Amount struct {
	decimal.Decimal
}

var Zero = AmountFromInt(0)

func AmountFromDoge(amount float64) Amount {
	return Amount{
		decimal.NewFromFloat(amount),
	}
}

func AmountFromUint64(amt uint64) Amount {
	i := new(big.Int).SetUint64(amt)
	return Amount{
		decimal.NewFromBigInt(i, -8),
	}
}

func AmountFromInt64(amt int64) Amount {
	return Amount{decimal.NewFromInt(amt)}
}

func AmountFromInt(amt int) Amount {
	return AmountFromInt64(int64(amt))
}

func AmountFromDecimal(amount decimal.Decimal) Amount {
	return Amount{
		amount,
	}
}

func (d Amount) ToDoge() float64 {
	amt, _ := d.Decimal.Truncate(8).Float64()
	return amt
}

func (d Amount) ToUint64() uint64 {
	amt := d.Decimal.Mul(decimal.NewFromInt(1e8)).IntPart()
	return uint64(amt)
}

func (d Amount) Add(d2 Amount) Amount {
	return Amount{
		Decimal: d.Decimal.Add(d2.Decimal),
	}
}

func (d Amount) Sub(d2 Amount) Amount {
	return Amount{
		Decimal: d.Decimal.Sub(d2.Decimal),
	}
}

func (d Amount) Div(d2 Amount) Amount {
	return Amount{
		Decimal: d.Decimal.Div(d2.Decimal),
	}
}

func (d Amount) Mul(d2 Amount) Amount {
	return Amount{
		Decimal: d.Decimal.Mul(d2.Decimal),
	}
}

func (d Amount) LessThan(d2 Amount) bool {
	return d.Decimal.LessThan(d2.Decimal)
}

func (d Amount) GreaterThan(d2 Amount) bool {
	return d.Decimal.GreaterThan(d2.Decimal)
}

func (d Amount) GreaterThanOrEqual(d2 Amount) bool {
	return d.Decimal.GreaterThanOrEqual(d2.Decimal)
}
