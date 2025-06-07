package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/db"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewCardanoHdWallet(config *domain.ADAConfig, logger logger.Logger) domain.CardanoWallet {
	hd := new(cardanoWallet)
	hd.MasterKey = createMasterKey(config, logger)
	hd.MasterWallet = db.CreateMasterWallet(config, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterWallet, logger)
	return hd
}
