package domain

import (
	"encoding/hex"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/btcsuite/btcutil/base58"
	"github.com/islishude/bip32"
)

type NEARConfig struct {
	config.SecureConfig
}

type MasterKey struct {
	Mnemonic  string
	Seed      []byte
	MasterKey bip32.XPrv
}

type WalletMainAcc struct {
	PrivateKey *Ed25519KeyPair
}

type NearWallet interface {
	domain.AccountManagerModel
	AccPrv(path ad.DerivationPath) (*Ed25519KeyPair, error)
}

func (k *MasterKey) GenerateEd25519KeyPair(path ad.DerivationPath) (*Ed25519KeyPair, error) {
	var (
		kp Ed25519KeyPair
	)

	priv := k.MasterKey.DeriveHard(path.Account).Derive(path.Index)

	kp.Ed25519PrivKey = priv
	kp.Ed25519PubKey = priv.XPub()
	kp.PublicKey = ed25519Prefix + base58.Encode(kp.Ed25519PubKey.PublicKey())
	kp.PrivateKey = ed25519Prefix + base58.Encode(kp.Ed25519PrivKey.Bytes())
	kp.AccountID = hex.EncodeToString(kp.Ed25519PubKey.PublicKey())
	return &kp, nil
}
