package domain

import (
	"crypto/ecdsa"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
)

type TRXConfig struct {
	config.SecureConfig
	ChainType string
	NetworkId int64
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	Secret    [32]byte
	ChainCode [32]byte
	MasterKey *hdkeychain.ExtendedKey
}

type WalletMainAcc struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Account    *keystore.Account
	Path       string
}

type TronWallet interface {
	domain.AccountManagerModel
	MasterAccPrv() *ecdsa.PrivateKey
	AccPrv(path string) (*ecdsa.PrivateKey, error)
}
