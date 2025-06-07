package wallet

import (
	"crypto/ecdsa"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/matic/domain"
	"github.com/ethereum/go-ethereum/accounts"
)

type maticWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w maticWallet) Coin() string {
	return "MATIC"
}

func (w maticWallet) GetAccount(request ad.Request) (*ad.Account, error) {
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

func (w maticWallet) AccPrv(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	prv, err := privateFromMasterAndAddressBTC(w.MasterKey.MasterKey, path)
	if err != nil {
		return nil, err
	}
	return prv.ToECDSA(), nil
}
