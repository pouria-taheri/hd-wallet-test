package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/islishude/bip32"
)

func createMasterKey(config config.SecureConfig, log logger.Logger) domain.MasterKey {
	ma := domain.MasterKey{}
	seed, err := config.GetSeedWithErrorChecking()
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create seed from mnemonic")
	}
	ma.Seed = seed
	ma.MasterKey = bip32.NewRootXPrv(seed)

	return ma
}
