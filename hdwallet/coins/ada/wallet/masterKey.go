package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func createMasterKey(config *domain.ADAConfig, log logger.Logger) domain.MasterKey {

	ma := domain.MasterKey{}
	seed, err := config.GetSeedWithErrorChecking()
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create seed from mnemonic")
	}
	ma.Seed = seed

	return ma
}
