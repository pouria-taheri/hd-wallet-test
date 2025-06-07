package wallet

import (
	"crypto/ecdsa"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/avax/domain"
	"github.com/ethereum/go-ethereum/accounts"
)

type avaxWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w avaxWallet) Coin() string {
	return "AVAX"
}

func (w avaxWallet) GetAccount(request ad.Request) (*ad.Account, error) {
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

func (w avaxWallet) AccPrv(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	prv, err := privateFromMasterAndAddressBTC(w.MasterKey.MasterKey, path)
	if err != nil {
		return nil, err
	}
	return prv.ToECDSA(), nil
}
