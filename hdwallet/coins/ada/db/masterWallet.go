package db

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/address"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

type MasterWallet struct {
	Wallet *Wallet
	DB     DB
}

func CreateMasterWallet(config *domain.ADAConfig, log logger.Logger) MasterWallet {
	mw := MasterWallet{}

	entropy, err := bip39.EntropyFromMnemonic(config.Mnemonic)
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterWallet",
			"error":     err,
		}).FatalF("cannot create entropy from mnemonic")
	}

	mw.DB = NewBadgerDB(*config)

	wallet, err := mw.DB.GetWallet(config.WalletName)
	if err != nil {
		if err == WalletNotFoundErr {
			wallet = NewWallet(config.WalletName, config.Password, entropy)
			err = mw.DB.SaveWallet(wallet)
			if err != nil {
				log.With(logger.Field{
					"submodule": "masterWallet",
					"error":     err,
				}).FatalF("cannot save wallet into db")
			}
		} else {
			log.With(logger.Field{
				"submodule": "masterWallet",
				"error":     err,
			}).FatalF("cannot get wallet")
		}
	}

	if wallet == nil {
		log.With(logger.Field{
			"submodule": "masterWallet",
			"error":     err,
		}).FatalF("wallet is empty")
	}

	wallet.SetNetwork(getNetwork(config.ChainType))
	mw.Wallet = wallet

	return mw
}

func getNetwork(chainType string) (network address.Network) {
	switch strings.ToLower(chainType) {
	case "mainnet":
		network = address.Mainnet
	case "testnet":
		network = address.Testnet
	default:
		network = address.Testnet
	}
	return network
}
