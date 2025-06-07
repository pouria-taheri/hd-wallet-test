package usecase

import (
	"encoding/hex"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/btcd/txscript"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"sync"
)

// managedAddress represents a public key address.  It also may or may not have
// the private key associated with the public key.
type managedAddress struct {
	manager        *ScopedKeyManager
	derivationPath ad.DerivationPath
	address        btcutil.Address
	imported       bool
	internal       bool
	compressed     bool
	used           bool
	addrType       ad.AddressType
	pubKey         *btcec.PublicKey
	prvKey         []byte
	privKeyCT      []byte // non-nil if unlocked
	privKeyMutex   sync.Mutex
}

// Account returns the account number the address is associated with.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) Account() uint32 {
	return a.derivationPath.Account
}

// AddrType returns the address type of the managed address. This can be used
// to quickly discern the address type without further processing
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) AddrType() ad.AddressType {
	return a.addrType
}

// Address returns the btcutil.Address which represents the managed address.
// This will be a pay-to-pubkey-hash address.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) Address() ad.Address {
	return a.address
}

// AddrHash returns the public key hash for the address.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) AddrHash() []byte {
	var hash []byte

	switch n := a.address.(type) {
	case *btcutil.AddressPubKeyHash:
		hash = n.Hash160()[:]
	case *btcutil.AddressScriptHash:
		hash = n.Hash160()[:]
	case *btcutil.AddressWitnessPubKeyHash:
		hash = n.Hash160()[:]
	}

	return hash
}

// Imported returns true if the address was imported instead of being part of an
// address chain.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) Imported() bool {
	return a.imported
}

// Internal returns true if the address was created for internal use such as a
// change output of a transaction.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) Internal() bool {
	return a.internal
}

// Compressed returns true if the address is compressed.
//
// This is part of the ManagedAddress interface implementation.
func (a *managedAddress) Compressed() bool {
	return a.compressed
}

// PubKey returns the public key associated with the address.
//
// This is part of the ManagedPubKeyAddress interface implementation.
func (a *managedAddress) PubKey() *btcec.PublicKey {
	return a.pubKey
}

// pubKeyBytes returns the serialized public key bytes for the managed address
// based on whether or not the managed address is marked as compressed.
func (a *managedAddress) pubKeyBytes() []byte {
	if a.compressed {
		return a.pubKey.SerializeCompressed()
	}
	return a.pubKey.SerializeUncompressed()
}

// ExportPubKey returns the public key associated with the address
// serialized as a hex encoded string.
//
// This is part of the ManagedPubKeyAddress interface implementation.
func (a *managedAddress) ExportPubKey() string {
	return hex.EncodeToString(a.pubKeyBytes())
}

// PrivKey returns the private key for the address.  It can fail if the address
// manager is watching-only or locked, or the address does not have any keys.
//
// This is part of the ManagedPubKeyAddress interface implementation.
func (a *managedAddress) PrivKey() (*btcec.PrivateKey, error) {
	a.manager.mtx.Lock()
	defer a.manager.mtx.Unlock()

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), a.prvKey)
	return privKey, nil
}

// ExportPrivKey returns the private key associated with the address in Wallet
// Import Format (WIF).
//
// This is part of the ManagedPubKeyAddress interface implementation.
func (a *managedAddress) ExportPrivKey() (*btcutil.WIF, error) {
	pk, err := a.PrivKey()
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(pk, a.manager.rootManager.chainParams, a.compressed)
}

// DerivationInfo contains the information required to derive the key that
// backs the address via traditional methods from the HD root. For imported
// keys, the first value will be set to false to indicate that we don't know
// exactly how the key was derived.
//
// This is part of the ManagedPubKeyAddress interface implementation.
func (a *managedAddress) DerivationInfo() (ad.KeyScope, ad.DerivationPath, bool) {
	var (
		scope ad.KeyScope
		path  ad.DerivationPath
	)

	// If this key is imported, then we can't return any information as we
	// don't know precisely how the key was derived.
	if a.imported {
		return scope, path, false
	}

	return a.manager.Scope(), a.derivationPath, true
}

// newManagedAddressWithoutPrivKey returns a new managed address based on the
// passed account, public key, and whether or not the public key should be
// compressed.
func newManagedAddressWithoutPrivKey(m *ScopedKeyManager, derivationPath ad.DerivationPath,
	pubKey *btcec.PublicKey, compressed bool, addrType ad.AddressType) (*managedAddress, error) {

	// Create a pay-to-pubkey-hash address from the public key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	var address btcutil.Address
	var err error

	switch addrType {

	case ad.NestedWitnessPubKey:
		// For this address type we'l generate an address which is
		// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
		// allows us to take advantage of segwit's scripting improvments,
		// and malleability fixes.

		// First, we'll generate a normal p2wkh address from the pubkey hash.
		witAddr, err := btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, m.rootManager.chainParams,
		)
		if err != nil {
			return nil, err
		}

		// Next we'll generate the witness program which can be used as a
		// pkScript to pay to this generated address.
		witnessProgram, err := txscript.PayToAddrScript(witAddr)
		if err != nil {
			return nil, err
		}

		// Finally, we'll use the witness program itself as the pre-image
		// to a p2sh address. In order to spend, we first use the
		// witnessProgram as the sigScript, then present the proper
		// <sig, pubkey> pair as the witness.
		address, err = btcutil.NewAddressScriptHash(
			witnessProgram, m.rootManager.chainParams,
		)
		if err != nil {
			return nil, err
		}

	case ad.PubKeyHash:
		address, err = btcutil.NewAddressPubKeyHash(
			pubKeyHash, m.rootManager.chainParams,
		)
		if err != nil {
			return nil, err
		}

	case ad.WitnessPubKey:
		address, err = btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, m.rootManager.chainParams,
		)
		if err != nil {
			return nil, err
		}
	}

	return &managedAddress{
		manager:        m,
		address:        address,
		derivationPath: derivationPath,
		imported:       false,
		internal:       false,
		addrType:       addrType,
		compressed:     compressed,
		pubKey:         pubKey,
		prvKey:         nil,
		privKeyCT:      nil,
	}, nil
}

// newManagedAddress returns a new managed address based on the passed account,
// private key, and whether or not the public key is compressed.  The managed
// address will have access to the private and public keys.
func newManagedAddress(s *ScopedKeyManager, derivationPath ad.DerivationPath,
	privKey *btcec.PrivateKey, compressed bool, addrType ad.AddressType) (*managedAddress, error) {

	// Encrypt the private key.
	//
	// NOTE: The privKeyBytes here are set into the managed address which
	// are cleared when locked, so they aren't cleared here.
	privKeyBytes := privKey.Serialize()

	// Leverage the code to create a managed address without a private key
	// and then add the private key to it.
	ecPubKey := (*btcec.PublicKey)(&privKey.PublicKey)
	managedAddr, err := newManagedAddressWithoutPrivKey(
		s, derivationPath, ecPubKey, compressed, addrType,
	)
	if err != nil {
		return nil, err
	}
	managedAddr.prvKey = privKeyBytes
	managedAddr.privKeyCT = privKeyBytes

	return managedAddr, nil
}

// newManagedAddressFromExtKey returns a new managed address based on the passed
// account and extended key.  The managed address will have access to the
// private and public keys if the provided extended key is private, otherwise it
// will only have access to the public key.
func newManagedAddressFromExtKey(s *ScopedKeyManager, derivationPath ad.DerivationPath,
	key *hdkeychain.ExtendedKey, addrType ad.AddressType) (*managedAddress, error) {

	// Create a new managed address based on the public or private key
	// depending on whether the generated key is private.
	var managedAddr *managedAddress
	if key.IsPrivate() {

		privKey, err := key.ECPrivKey()
		if err != nil {
			s.logger.With(logger.Field{
				"submodule": "address usecase",
				"section":   "manageAddressFromExtendKey",
				"error":     err,
			}).ErrorF("cannot get private key from extended key")
			return nil, err
		}

		// Ensure the temp private key big integer is cleared after
		// use.
		managedAddr, err = newManagedAddress(
			s, derivationPath, privKey, true, addrType,
		)
		if err != nil {
			s.logger.With(logger.Field{
				"submodule": "address usecase",
				"section":   "manageAddressFromExtendKey",
				"error":     err,
				"scope":     s,
				"dPath":     derivationPath,
				"addrType":  addrType,
			}).ErrorF("cannot get managed address from private key")
			return nil, err
		}
	} else {
		pubKey, err := key.ECPubKey()
		if err != nil {
			s.logger.With(logger.Field{
				"submodule": "address usecase",
				"section":   "manageAddressFromExtendKey",
				"error":     err,
			}).ErrorF("cannot get public key from extended key")
			return nil, err
		}

		managedAddr, err = newManagedAddressWithoutPrivKey(
			s, derivationPath, pubKey, true,
			addrType,
		)
		if err != nil {
			s.logger.With(logger.Field{
				"submodule": "address usecase",
				"section":   "manageAddressFromExtendKey",
				"error":     err,
				"scope":     s,
				"dPath":     derivationPath,
				"addrType":  addrType,
			}).ErrorF("cannot get managed address from public key")
			return nil, err
		}
	}

	return managedAddr, nil
}
