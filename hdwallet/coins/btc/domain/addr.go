package domain

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type Address struct {
	DerivationPath
	KeyScope
	// Private returns either a public or private derived extended key
	// based on the flag state
	Private bool `json:"private"`

	Address string `json:"address"`

	BtcAddr btcutil.Address `json:"-"`
}

func (a *Address) GenerateAddress(chainParams *chaincfg.Params) error {
	addr, err := btcutil.DecodeAddress(a.Address, chainParams)
	if err != nil {
		return err
	}
	a.BtcAddr = addr
	return nil
}
