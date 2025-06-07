package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcutil"
)

type Address struct {
	ad.DerivationPath
	ad.KeyScope
	// Private returns either a public or private derived extended key
	// based on the flag state
	Private bool `json:"private"`

	Address string `json:"address"`

	LtcAddr ltcutil.Address `json:"-"`
}

func (a *Address) GenerateAddress(chainParams *chaincfg.Params) error {
	addr, err := ltcutil.DecodeAddress(a.Address, chainParams)
	if err != nil {
		return err
	}
	a.LtcAddr = addr
	return nil
}
