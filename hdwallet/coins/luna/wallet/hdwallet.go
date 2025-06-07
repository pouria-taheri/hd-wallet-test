package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewTerraHdWallet(config *domain.LUNAConfig, logger logger.Logger) domain.TerraWallet {
	hd := new(terraWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
