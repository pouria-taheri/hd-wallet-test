package wallet

import (
	"crypto/ecdsa"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/trx/domain"
)

type tronWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w tronWallet) Coin() string {
	return "TRX"
}

func (w tronWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	path := pathFromUserID(uint64(request.Account), request.Index)
	prv, _, err := privateFromMasterAndAddress(w.Secret, w.ChainCode, path)
	if err != nil {
		return nil, err
	}
	addr := addressFromPrivate(prv)
	return &ad.Account{
		Address: addr.Address.String(),
	}, nil
}

func (w tronWallet) AccPrv(path string) (*ecdsa.PrivateKey, error) {
	prv, _, err := privateFromMasterAndAddress(w.Secret, w.ChainCode, path)
	if err != nil {
		return nil, err
	}
	return prv.ToECDSA(), nil

}
