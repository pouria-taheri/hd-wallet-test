package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/db"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func createWalletMainAccount(key db.MasterWallet, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	wma.Account = key.Wallet.GetAddress(0, 0)
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.Account,
	}).InfoF("Main account")

	return wma
}
