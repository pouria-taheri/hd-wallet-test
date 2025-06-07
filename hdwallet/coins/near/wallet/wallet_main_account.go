package wallet

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/near/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	prv, err := key.GenerateEd25519KeyPair(ad.DerivationPath{Index: 0, Account: 0})
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     err,
		}).ErrorF("error in get master private key")
		panic("error in get master private key from path")
	}
	wma.PrivateKey = prv

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.PrivateKey.AccountID,
	}).InfoF("Main account")

	return wma
}
