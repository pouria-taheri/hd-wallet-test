package wallet

import (
	"context"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/eos/domain"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"github.com/eoscanada/eos-go/ecc"
)

type eosWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
	domain.WalletChildAcc
}

func (w *eosWallet) Coin() string {
	return "EOS"
}

func (w *eosWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 32765294217
	path := pathFromUserId(uint64(request.Account), request.Index)

	prv, err := w.AccPrv(path)
	if err != nil {
		return nil, err
	}

	mp, err := w.AccPrv(pathFromUserId(0, 0))
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Master:  mp.PublicKey().String(),
		Address: prv.PublicKey().String(),
		Memo:    fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
	}, nil
}

func (w *eosWallet) AccPrv(path string) (*ecc.PrivateKey, error) {
	key, errKey := domain.DeriveForPath(path, w.Seed)
	if errKey != nil {
		return nil, errKey
	}

	wif := base58.CheckEncode(key.Key, '\x80')
	w.WalletChildAcc.Wif = wif
	privateKey, errPrv := ecc.NewPrivateKey(wif)
	if errPrv != nil {
		return nil, errPrv
	}
	return privateKey, nil
}

func (w *eosWallet) AccKeyBag(path string) (*eos.KeyBag, error) {
	prv, err := w.AccPrv(path)
	if err != nil {
		return nil, err
	}
	keyBag := &eos.KeyBag{}
	errKey := keyBag.ImportPrivateKey(context.Background(), w.Wif)
	if errKey != nil {
		return nil, errKey
	}
	w.WalletChildAcc.PrivateKey = prv
	w.WalletChildAcc.PublicKey = prv.PublicKey()

	return keyBag, nil
}

func (w *eosWallet) AccPublicKey() ecc.PublicKey {
	return w.WalletChildAcc.PublicKey
}
