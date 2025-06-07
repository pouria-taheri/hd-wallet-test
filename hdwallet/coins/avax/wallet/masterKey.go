package wallet

import (
	"crypto/ecdsa"
	"git.mazdax.tech/blockchain/hdwallet/coins/avax/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func createMasterKey(config config.SecureConfig, net *chaincfg.Params, log logger.Logger) domain.MasterKey {
	ma := domain.MasterKey{}
	seed, err := config.GetSeedWithErrorChecking()
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).ErrorF("cannot create seed from mnemonic")
		panic("cannot create seed from mnemonic")
	}
	ma.Seed = seed
	mk, err := MasterKeyFromSeed(seed, net)
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).ErrorF("cannot create masterKey from seed")
		panic("cannot create masterKey from seed")
	}
	ma.MasterKey = mk
	return ma
}

func MasterKeyFromSeed(seed []byte, net *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	return hdkeychain.NewMaster(seed, net)
}

func (w avaxWallet) MasterAccPrv() *ecdsa.PrivateKey { return w.WalletMainAcc.PrivateKey }
