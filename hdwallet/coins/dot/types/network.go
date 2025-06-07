package types

import "fmt"

type Network uint8

const (
	Westend Network = iota
	Polkadot
	Kusama
	Substrate
)

func (n Network) String() string {
	switch n {
	case Westend:
		return "westend"
	case Polkadot:
		return "polkadot"
	case Kusama:
		return "kusama"
	case Substrate:
		return "substrate"
	default:
		return ""
	}
}

func (n Network) SS58Prefix() uint8 {
	switch n {
	case Polkadot:
		return 0
	case Kusama:
		return 2
	case Westend:
		return 42
	case Substrate:
		return 42
	default:
		return 0
	}
}

func (n Network) BIP44CoinType() string {
	switch n {
	case Polkadot:
		return fmt.Sprintf("%d", 0x80000162)
	case Kusama:
		return fmt.Sprintf("%d", 0x800001b2)
	case Westend:
		return "westend"
	case Substrate:
		return ""
	default:
		return ""
	}
}

func IsValidNetwork(str string) bool {
	return str == string(Polkadot) || str == string(Kusama) || str == string(Westend) || str == string(Substrate)
}
