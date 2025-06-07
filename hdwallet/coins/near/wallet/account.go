package wallet

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/near/domain"
)

type nearWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *nearWallet) Coin() string {
	return "NEAR"
}

func (w *nearWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	prv, err := w.AccPrv(ad.DerivationPath{Index: request.Index, Account: request.Account})
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Id:               uint64(request.Account),
		Address:          prv.AccountID,
		PublicKey:        []byte(prv.PublicKey),
		Ed25519PublicKey: prv.Ed25519PubKey.PublicKey(),
	}, nil
}

func (w *nearWallet) AccPrv(path ad.DerivationPath) (*domain.Ed25519KeyPair, error) {
	return w.MasterKey.GenerateEd25519KeyPair(path)
}
