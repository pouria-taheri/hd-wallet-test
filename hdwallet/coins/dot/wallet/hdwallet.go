package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewDotHdWallet(config *domain.DOTConfig, logger logger.Logger) domain.DotWallet {
	hd := new(dotWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, config.NetworkId, logger)
	hd.networkId = config.NetworkId
	return hd
}
