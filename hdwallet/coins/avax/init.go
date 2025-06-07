package avax

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/avax/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/avax/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) baseDomain.AccountManagerModel {
	c := new(domain.AvaxConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewAvaxHdWallet(c, logger)
	return hd
}
