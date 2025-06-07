package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewBinanceHdWallet(config *domain.BNBConfig, logger logger.Logger) domain.BinanceWallet {
	hd := new(binanceWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, config.NetworkId, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
