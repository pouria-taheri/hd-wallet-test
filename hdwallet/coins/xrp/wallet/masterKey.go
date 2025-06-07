package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
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
	key, errKey := domain.DeriveForPath(pathFromUserId(0, 0), seed)
	if errKey != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create derivation from seed")
		return ma
	}
	ma.MasterKey = key

	return ma
}
