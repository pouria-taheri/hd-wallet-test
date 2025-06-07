package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
)

type XRPConfig struct {
	config.SecureConfig
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey *Key
}

type WalletMainAcc struct {
	PrivateKey *Key
	PublicKey  string
}

type XrpWallet interface {
	domain.AccountManagerModel
	AccPrv(path string) (*Key, error)
}
