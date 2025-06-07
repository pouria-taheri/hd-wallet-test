package domain

import (
	"github.com/islishude/bip32"
)

const ed25519Prefix = "ed25519:"

// Ed25519KeyPair is a Ed25519 key pair.
type Ed25519KeyPair struct {
	AccountID      string     `json:"account_id"`
	PublicKey      string     `json:"public_key"`
	PrivateKey     string     `json:"private_key,omitempty"`
	Ed25519PubKey  bip32.XPub `json:"-"`
	Ed25519PrivKey bip32.XPrv `json:"-"`
}
