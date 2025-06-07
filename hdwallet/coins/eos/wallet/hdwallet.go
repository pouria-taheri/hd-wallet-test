package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/eos/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewEosHdWallet(config *domain.EOSConfig, logger logger.Logger) domain.EosWallet {
	hd := new(eosWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey, logger)
	return hd
}
