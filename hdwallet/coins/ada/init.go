package ada

import (
	domain2 "git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/wallet"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewCardanoWalletUseCase(config configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) domain2.CardanoWalletModel {
	c := new(domain.ADAConfig)
	if err := config.Unmarshal(c); err != nil {
		panic(err)
	}
	c.SecureConfig = secureConfig
	hd := wallet.NewCardanoHdWallet(c, logger)
	return hd
}
