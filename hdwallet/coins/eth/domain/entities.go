package domain

import (
	"crypto/ecdsa"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

type EthConfig struct {
	config.SecureConfig
	ChainType string
	NetworkId int64
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey *hdkeychain.ExtendedKey
}

type WalletMainAcc struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    *common.Address
	Path       accounts.DerivationPath
}

type ETHWallet interface {
	domain.AccountManagerModel
	MasterAccPrv() *ecdsa.PrivateKey
	AccPrv(path accounts.DerivationPath) (*ecdsa.PrivateKey, error)
}
