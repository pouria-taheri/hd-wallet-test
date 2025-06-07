package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"github.com/ltcsuite/ltcd/btcec"
	"github.com/ltcsuite/ltcutil"
)

// ManagedPubKeyAddress extends ManagedAddress and additionally provides the
// public and private keys for pubkey-based addresses.
type ManagedPubKeyAddress interface {
	ad.ManagedAddress

	// PubKey returns the public key associated with the address.
	PubKey() *btcec.PublicKey

	// ExportPubKey returns the public key associated with the address
	// serialized as a hex encoded string.
	ExportPubKey() string

	// PrivKey returns the private key for the address.  It can fail if the
	// address manager is watching-only or locked, or the address does not
	// have any keys.
	PrivKey() (*btcec.PrivateKey, error)

	// ExportPrivKey returns the private key associated with the address
	// serialized as Wallet Import Format (WIF).
	ExportPrivKey() (*ltcutil.WIF, error)

	// DerivationInfo contains the information required to derive the key
	// that backs the address via traditional methods from the HD root. For
	// imported keys, the first value will be set to false to indicate that
	// we don't know exactly how the key was derived.
	DerivationInfo() (ad.KeyScope, ad.DerivationPath, bool)
}

// ManagedScriptAddress extends ManagedAddress and represents a pay-to-script-hash
// style of bitcoin addresses.  It additionally provides information about the
// script.
type ManagedScriptAddress interface {
	ad.ManagedAddress

	// Script returns the script associated with the address.
	Script() ([]byte, error)
}
