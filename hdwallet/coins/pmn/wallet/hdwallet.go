package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/pmn/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewKuknosHdWallet(config *domain.PMNConfig, logger logger.Logger) domain.KuknosWallet {
	hd := new(stellarWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
