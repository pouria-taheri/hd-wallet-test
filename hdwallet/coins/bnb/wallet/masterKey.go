package wallet

import (
	"git.mazdax.tech/blockchain/bnb-go-sdk/common/types"
	"git.mazdax.tech/blockchain/bnb-go-sdk/keys"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func createMasterKey(config config.SecureConfig, networkId types.ChainNetwork, log logger.Logger) domain.MasterKey {
	seed, err := config.GetSeedWithErrorChecking()
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).ErrorF("cannot create seed from mnemonic")
		panic("cannot create seed from mnemonic")
	}

	ma := domain.MasterKey{}
	ma.Seed = seed
	master, ch := MasterKeyFromSeed(seed)
	ma.ChainCode = ch
	ma.Secret = master
	ma.Mnemonic = config.Mnemonic

	types.Network = networkId
	path := pathFromUserID(0, 0)
	ma.KeyManager, err = keys.NewKeyManagerWithSecretChainCode(ma.Secret, ma.ChainCode, path)
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"path":      path,
			"error":     err,
		}).ErrorF("cannot create key manager from secret key and path")
		panic("cannot create key manager")
	}

	ma.PrivKey = ma.KeyManager.GetPrivKey()

	return ma
}

func MasterKeyFromSeed(seed []byte) (secret, chainCode [32]byte) {
	secret, chainCode = keys.ComputeMastersFromSeed(seed)
	return
}
