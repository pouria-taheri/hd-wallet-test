package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/xlm/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewStellarHdWallet(config *domain.XLMConfig, logger logger.Logger) domain.StellarWallet {
	hd := new(stellarWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
