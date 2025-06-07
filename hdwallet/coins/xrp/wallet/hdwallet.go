package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewXrpHdWallet(config *domain.XRPConfig, logger logger.Logger) domain.XrpWallet {
	hd := new(xrpWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
