package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/helper"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/types"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/amintalebi/go-subkey"
	"github.com/amintalebi/go-subkey/ed25519"
)

type DOTConfig struct {
	config.SecureConfig
	NetworkId types.Network
}

type MasterKey struct {
	Mnemonic string
	Seed     []byte
	Key      subkey.KeyPair
	Scheme   ed25519.Scheme
}

type WalletMainAcc struct {
	PrivateKey subkey.KeyPair
}

type DotWallet interface {
	domain.AccountManagerModel
	AccPrv(path helper.DerivationPath) (subkey.KeyPair, error)
}
