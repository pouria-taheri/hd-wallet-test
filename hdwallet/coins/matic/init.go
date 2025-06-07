package matic

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/matic/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/matic/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) baseDomain.AccountManagerModel {
	c := new(domain.MaticConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewMaticHdWallet(c, logger)
	return hd
}
