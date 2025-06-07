package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/bnb-go-sdk/keys"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type binanceWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *binanceWallet) Coin() string {
	return "BNB"
}

func (w *binanceWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 24527364358
	path := pathFromUserID(uint64(request.Account), 0)

	key, err := keys.NewKeyManagerWithSecretChainCode(w.Secret, w.ChainCode, path)
	if err != nil {
		return nil, err
	}

	mk, err := keys.NewKeyManagerWithSecretChainCode(w.Secret, w.ChainCode, pathFromUserID(0, 0))
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Address: key.GetAddr().String(),
		Memo:    fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
		Master:  mk.GetAddr().String(),
	}, nil
}

func (w *binanceWallet) AccPrv(path string) (crypto.PrivKey, error) {
	derivedPriv, err := keys.DerivePrivateKeyForPath(w.Secret, w.ChainCode, path)
	if err != nil {
		return nil, err
	}
	return secp256k1.GenPrivKeySecp256k1(derivedPriv), nil
}
