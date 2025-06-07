package wallet

import (
	"crypto/ecdsa"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ftm/domain"
	"github.com/ethereum/go-ethereum/accounts"
)

type ftmWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w ftmWallet) Coin() string {
	return "FTM"
}

func (w ftmWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	path, err := pathFromUserID(uint64(request.Account), request.Index)
	if err != nil {
		return nil, err
	}
	prv, err := privateFromMasterAndAddressBTC(w.MasterKey.MasterKey, path)
	if err != nil {
		return nil, err
	}
	addr := addressFromPrivate(prv)
	return &ad.Account{
		Address: addr.String(),
	}, nil
}

func (w ftmWallet) AccPrv(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	prv, err := privateFromMasterAndAddressBTC(w.MasterKey.MasterKey, path)
	if err != nil {
		return nil, err
	}
	return prv.ToECDSA(), nil
}
