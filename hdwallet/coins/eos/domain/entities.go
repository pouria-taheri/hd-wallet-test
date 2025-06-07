package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type EOSConfig struct {
	config.SecureConfig
	MasterAccountName string
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey *Key
}

type WalletMainAcc struct {
	PrivateKey *ecc.PrivateKey
	PublicKey  ecc.PublicKey
}

type WalletChildAcc struct {
	PrivateKey *ecc.PrivateKey
	PublicKey  ecc.PublicKey
	Wif        string
}

type EosWallet interface {
	domain.AccountManagerModel
	AccPrv(path string) (*ecc.PrivateKey, error)
	AccKeyBag(path string) (*eos.KeyBag, error)
	AccPublicKey() ecc.PublicKey
}
