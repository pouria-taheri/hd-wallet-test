package domain

import (
	"git.mazdax.tech/blockchain/bnb-go-sdk/common/types"
	"git.mazdax.tech/blockchain/bnb-go-sdk/keys"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/tendermint/tendermint/crypto"
)

type BNBConfig struct {
	config.SecureConfig
	NetworkId types.ChainNetwork
}

type MasterKey struct {
	Mnemonic   string
	KeyManager keys.KeyManager
	Seed       []byte
	Secret     [32]byte
	ChainCode  [32]byte
	PrivKey    crypto.PrivKey
}

type WalletMainAcc struct {
	Account types.AccAddress
	Path    string
}

type BinanceWallet interface {
	domain.AccountManagerModel
	MasterAccPrv() crypto.PrivKey
	AccPrv(path string) (crypto.PrivKey, error)
	AccKeyManager(logger logger.Logger, account ad.Account) keys.KeyManager
}
