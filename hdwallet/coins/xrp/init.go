package xrp

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) baseDomain.AccountManagerModel {
	c := new(domain.XRPConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewXrpHdWallet(c, logger)
	return hd
}
