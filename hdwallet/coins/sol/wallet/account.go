package wallet

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	"github.com/gagliardetto/solana-go"
	"github.com/islishude/bip32"
)

type solWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *solWallet) Coin() string {
	return "SOL"
}

func (w *solWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	prv, err := w.AccPrv(ad.DerivationPath{Index: request.Index, Account: request.Account})
	if err != nil {
		return nil, err
	}

	public := solana.PublicKeyFromBytes(prv.PublicKey())

	return &ad.Account{
		Address: public.String(),
	}, nil
}

func (w *solWallet) AccPrv(path ad.DerivationPath) (bip32.XPrv, error) {
	return w.MasterKey.MasterKey.Derive(144).Derive(path.Account).Derive(path.Index), nil
}
