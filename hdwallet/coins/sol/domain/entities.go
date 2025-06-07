package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/gagliardetto/solana-go"
	"github.com/islishude/bip32"
)

type SOLConfig struct {
	config.SecureConfig
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey bip32.XPrv
}

type WalletMainAcc struct {
	PrivateKey bip32.XPrv
	PublicKey  solana.PublicKey
}

type SolWallet interface {
	domain.AccountManagerModel
	AccPrv(path ad.DerivationPath) (bip32.XPrv, error)
}
