package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewSolHdWallet(config *domain.SOLConfig, logger logger.Logger) domain.SolWallet {
	hd := new(solWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
