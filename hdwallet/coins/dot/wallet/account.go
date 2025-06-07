package wallet

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/helper"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/types"
	"github.com/amintalebi/go-subkey"
)

type dotWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
	networkId types.Network
}

func (w *dotWallet) Coin() string {
	return "DOT"
}

func (w *dotWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	path := w.dPathFromUserId(uint64(request.Account), request.Index)

	prv, err := w.AccPrv(path)
	if err != nil {
		return nil, err
	}

	address, err := prv.SS58Address(w.networkId.SS58Prefix())
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Address: address,
	}, nil
}

func (w *dotWallet) AccPrv(path helper.DerivationPath) (subkey.KeyPair, error) {
	djs, err := subkey.DeriveJunctions(subkey.DerivePath(path.String()))

	kp, err := w.Scheme.Derive(w.Key, djs)
	if err != nil {
		return nil, fmt.Errorf("error deriving from master. %w", err)
	}

	return kp, nil
}

func (w *dotWallet) dPathFromUserId(uid uint64, index uint32) helper.DerivationPath {
	return helper.DerivationPath{
		Network: w.networkId,
		Account: uid,
		Index:   index,
	}
}
