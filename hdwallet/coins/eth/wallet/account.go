package wallet

import (
	"crypto/ecdsa"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/eth/domain"
	"github.com/ethereum/go-ethereum/accounts"
)

type eTHWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w eTHWallet) Coin() string {
	return "ETH"
}

func (w eTHWallet) GetAccount(request ad.Request) (*ad.Account, error) {
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

func (w eTHWallet) AccPrv(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	prv, err := privateFromMasterAndAddressBTC(w.MasterKey.MasterKey, path)
	if err != nil {
		return nil, err
	}
	return prv.ToECDSA(), nil

}

//func (wmk MasterKey) WalletAddressByUserID(userID uint64) (string, error) {
//	dPath, err := eth.PathFromUserID(userID)
//	if err != nil {
//		return "", err
//	}
//	fmt.Println("dPath: " + dPath.String())
//	prv, err := privateFromMasterAndAddressBTC(wmk.MasterKey, dPath)
//	if err != nil {
//		return "", nil
//	}
//	addr := eth.AddressFromPrivate(prv)
//	return addr.String(), nil
//}
