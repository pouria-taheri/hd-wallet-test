package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/atom/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewCosmosHdWallet(config *domain.AtomConfig, logger logger.Logger) domain.CosmosWallet {
	hd := new(cosmosWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
