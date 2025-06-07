package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

type AtomConfig struct {
	config.SecureConfig
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	Secret    [32]byte
	ChainCode [32]byte
}

type WalletMainAcc struct {
	PrivateKey []byte
}

type CosmosWallet interface {
	domain.AccountManagerModel
	AccPrv(path *hd.BIP44Params) (tmsecp256k1.PrivKey, error)
}
