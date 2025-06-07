package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/aliworkshop/stellar-go/exp/crypto/derivation"
	"github.com/aliworkshop/stellar-go/keypair"
)

type XLMConfig struct {
	config.SecureConfig
	ChainType string
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey *derivation.Key
}

type WalletMainAcc struct {
	PublicKey []byte
	Account   *keypair.Full
}

type StellarWallet interface {
	domain.AccountManagerModel
	AccPrv(path string) (*derivation.Key, error)
}
