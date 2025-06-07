package domain

import (
	domain2 "git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/address"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/crypto"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
)

type ADAConfig struct {
	config.SecureConfig
	ChainType  string
	WalletName string
	DBPath     string
	BackupPath string
}

type MasterKey struct {
	Mnemonic string
	Seed     []byte
}

type WalletMainAcc struct {
	Account address.Address
}

type CardanoWallet interface {
	domain.AccountManagerModel
	domain2.CardanoWalletModel
	AccPrv(accountId, index uint32) (crypto.ExtendedSigningKey, error)
}
