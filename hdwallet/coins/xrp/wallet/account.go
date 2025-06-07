package wallet

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/crypto"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
)

type xrpWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *xrpWallet) Coin() string {
	return "XRP"
}

func (w *xrpWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 493517865
	path := pathFromUserId(uint64(request.Account), request.Index)

	prv, err := w.AccPrv(path)
	if err != nil {
		return nil, err
	}
	public, errPub := prv.PublicKey()
	if errPub != nil {
		return nil, errPub
	}

	xrpKey, _ := crypto.NewEd25519KeyFromPrivate(append(prv.Key, public...))
	address, _ := crypto.NewAccountId(xrpKey.Id(nil))

	return &ad.Account{
		Master:  w.PublicKey,
		Address: address.String(),
		Memo:    fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
	}, nil
}

func (w *xrpWallet) AccPrv(path string) (*domain.Key, error) {
	key, errKey := domain.DeriveForPath(path, w.Seed)
	if errKey != nil {
		return nil, errKey
	}

	return key, nil
}
