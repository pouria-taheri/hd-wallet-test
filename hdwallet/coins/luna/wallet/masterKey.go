package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
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

	master, ch := hd.ComputeMastersFromSeed(seed)
	ma.ChainCode = ch
	ma.Secret = master

	return ma
}
