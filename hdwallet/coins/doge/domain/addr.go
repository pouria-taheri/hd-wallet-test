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

	DogeAddr btcutil.Address `json:"-"`
}

func (a *Address) EncodeAddress() string {
	return a.DogeAddr.EncodeAddress()
}

func (a *Address) ScriptAddress() []byte {
	return a.DogeAddr.ScriptAddress()
}

func (a *Address) IsForNet(params *chaincfg.Params) bool {
	return a.DogeAddr.IsForNet(params)
}

func (a *Address) GetDerivationPath() DerivationPath {
	return a.DerivationPath
}

func (a *Address) GetKeyScope() KeyScope {
	return a.KeyScope
}

func (a *Address) IsPrivate() bool {
	return a.Private
}

func (a *Address) GenerateAddress(chainParams *chaincfg.Params) error {
	addr, err := btcutil.DecodeAddress(a.Address, chainParams)
	if err != nil {
		return err
	}
	a.DogeAddr = addr
	return nil
}
