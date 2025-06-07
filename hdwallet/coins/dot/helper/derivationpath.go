package helper

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/types"
)

type DerivationPath struct {
	Network types.Network
	Account uint64
	Index   uint32
}

func (d DerivationPath) String() string {
	return fmt.Sprintf("//m//44'//%s//%d//0//%d", d.Network.BIP44CoinType(), d.Account, d.Index)
}
