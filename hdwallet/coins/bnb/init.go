package bnb

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) domain.BinanceWallet {
	c := new(domain.BNBConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewBinanceHdWallet(c, logger)
	return hd
}
