package wallet

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/pmn/domain"
	"github.com/aliworkshop/stellar-go/exp/crypto/derivation"
	"github.com/aliworkshop/stellar-go/keypair"
)

type stellarWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *stellarWallet) Coin() string {
	return "PMN"
}

func (w *stellarWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 42540984621
	path := pathFromUserID(request.Account)
	prv, err := derivation.DeriveForPath(path, w.Seed)
	if err != nil {
		return nil, err
	}

	keyPair, err := keypair.FromRawSeed(prv.RawSeed())
	if err != nil {
		return nil, err
	}

	mp, err := derivation.DeriveForPath(pathFromUserID(0), w.Seed)
	if err != nil {
		return nil, err
	}

	mkp, err := keypair.FromRawSeed(mp.RawSeed())
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Master:  mkp.Address(),
		Address: keyPair.Address(),
		Memo:    fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
	}, nil
}

func (w *stellarWallet) AccPrv(path string) (*derivation.Key, error) {
	return derivation.DeriveForPath(path, w.Seed)
}
