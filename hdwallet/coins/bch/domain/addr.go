package domain

import (
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchutil"
)

type AddressModel interface {
	bchutil.Address
	GetDerivationPath() DerivationPath
	GetKeyScope() KeyScope
	IsPrivate() bool
}

type Address struct {
	DerivationPath
	KeyScope
	// Private returns either a public or private derived extended key
	// based on the flag state
	Private bool `json:"private"`

	Address string `json:"address"`

	BchAddr bchutil.Address `json:"-"`
}

func (a *Address) EncodeAddress() string {
	return a.BchAddr.EncodeAddress()
}

func (a *Address) ScriptAddress() []byte {
	return a.BchAddr.ScriptAddress()
}

func (a *Address) IsForNet(params *chaincfg.Params) bool {
	return a.BchAddr.IsForNet(params)
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
	addr, err := bchutil.DecodeAddress(a.Address, chainParams)
	if err != nil {
		return err
	}
	a.BchAddr = addr
	return nil
}
