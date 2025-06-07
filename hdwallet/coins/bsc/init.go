package bsc

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/bsc/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bsc/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) baseDomain.AccountManagerModel {
	c := new(domain.BscConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewBSCHdWallet(c, logger)
	return hd
}
