package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/pmn/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/aliworkshop/stellar-go/exp/crypto/derivation"
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
	mk, err := derivation.DeriveForPath("m/44'/1'/0'", seed)
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create masterKey from seed")
	}
	ma.MasterKey = mk
	return ma
}
