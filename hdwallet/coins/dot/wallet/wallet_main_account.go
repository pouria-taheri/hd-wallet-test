package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/helper"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/types"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/amintalebi/go-subkey"
)

func createWalletMainAccount(key domain.MasterKey, networkId types.Network, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}
	dPath := helper.DerivationPath{
		Network: networkId,
		Account: 0,
		Index:   0,
	}
	djs, err := subkey.DeriveJunctions(subkey.DerivePath(dPath.String()))
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create derive junction from dPath")
	}

	kp, err := key.Scheme.Derive(key.Key, djs)
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create key pair from master key")
	}
	wma.PrivateKey = kp

	address, err := kp.SS58Address(networkId.SS58Prefix())
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create address from key pair")
	}

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   address,
	}).InfoF("Main account")

	return wma
}
