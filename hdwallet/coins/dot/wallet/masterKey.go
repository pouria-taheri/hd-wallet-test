package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/amintalebi/go-subkey/ed25519"
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
	ma.Scheme = ed25519.Scheme{}
	kp, err := ma.Scheme.FromSeed(seed[:32])
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create key pair from seed")
		return ma
	}
	ma.Key = kp

	return ma
}
