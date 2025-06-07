package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/near/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewNearHdWallet(config *domain.NEARConfig, logger logger.Logger) domain.NearWallet {
	hd := new(nearWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
