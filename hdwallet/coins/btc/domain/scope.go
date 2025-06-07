package domain

import (
	"encoding/binary"
	"fmt"
)

// AddressType represents the various address types waddrmgr is currently able
// to generate, and maintain.
//
// NOTE: These MUST be stable as they're used for scope address schema
// recognition within the database.
type AddressType uint8

const (
	// PubKeyHash is a regular p2pkh address.
	PubKeyHash AddressType = iota

	// Script reprints a raw script address.
	Script

	// RawPubKey is just raw public key to be used within scripts, This
	// type indicates that a scoped manager with this address type
	// shouldn't be consulted during historical rescans.
	RawPubKey

	// NestedWitnessPubKey represents a p2wkh output nested within a p2sh
	// output. Using this address type, the wallet can receive funds from
	// other wallet's which don't yet recognize the new segwit standard
	// output types. Receiving funds to this address maintains the
	// scalability, and malleability fixes due to segwit in a backwards
	// compatible manner.
	NestedWitnessPubKey

	// WitnessPubKey represents a p2wkh (pay-to-witness-key-hash) address
	// type.
	WitnessPubKey
)

// DerivationPath represents a derivation path from a particular key manager's
// scope.  Each ScopedKeyManager starts key derivation from the end of their
// cointype hardened key: m/purpose'/cointype'. The fields in this struct allow
// further derivation to the next three child levels after the coin type key.
// This restriction is in the spriti of BIP0044 type derivation. We maintain a
// degree of coherency with the standard, but allow arbitrary derivations
// beyond the cointype key. The key derived using this path will be exactly:
// m/purpose'/cointype'/account/branch/index, where purpose' and cointype' are
// bound by the scope of a particular manager.
type DerivationPath struct {
	// Account is the account, or the first immediate child from the scoped
	// manager's hardened coin type key.
	Account uint32 `json:"account"`

	// Branch is the branch to be derived from the account index above. For
	// BIP0044-like derivation, this is either 0 (external) or 1
	// (internal). However, we allow this value to vary arbitrarily within
	// its size range.
	Branch uint32 `json:"branch"`

	// Index is the final child in the derivation path. This denotes the
	// key index within as a child of the account and branch.
	Index uint32 `json:"index"`
}

// KeyScope represents a restricted key scope from the primary root key within
// the HD chain. From the root manager (m/) we can create a nearly arbitrary
// number of ScopedKeyManagers of key derivation path: m/purpose'/cointype'.
// These scoped managers can then me managed indecently, as they house the
// encrypted cointype key and can derive any child keys from there on.
type KeyScope struct {
	// Purpose is the purpose of this key scope. This is the first child of
	// the master HD key.
	Purpose uint32 `json:"purpose"`

	// Coin is a value that represents the particular coin which is the
	// child of the purpose key. With this key, any accounts, or other
	// children can be derived at all.
	Coin uint32 `json:"coin"`
}

// scopeKeySize is the size of a scope as stored within the database.
const scopeKeySize = 8

// scopeToBytes transforms a manager's scope into the form that will be used to
// retrieve the bucket that all information for a particular scope is stored
// under
func (scope *KeyScope) Bytes() [scopeKeySize]byte {
	var scopeBytes [scopeKeySize]byte
	binary.LittleEndian.PutUint32(scopeBytes[:], scope.Purpose)
	binary.LittleEndian.PutUint32(scopeBytes[4:], scope.Coin)

	return scopeBytes
}

// String returns a human readable version describing the keypath encapsulated
// by the target key scope.
func (scope *KeyScope) String() string {
	return fmt.Sprintf("m/%v'/%v'", scope.Purpose, scope.Coin)
}

// ScopedIndex is a tuple of KeyScope and child Index. This is used to compactly
// identify a particular child key, when the account and branch can be inferred
// from context.
type ScopedIndex struct {
	// Scope is the BIP44 account' used to derive the child key.
	Scope KeyScope

	// Index is the BIP44 address_index used to derive the child key.
	Index uint32
}

// ScopeAddrSchema is the address schema of a particular KeyScope. This will be
// persisted within the database, and will be consulted when deriving any keys
// for a particular scope to know how to encode the public keys as addresses.
type ScopeAddrSchema struct {
	// ExternalAddrType is the address type for all keys within branch 0.
	ExternalAddrType AddressType

	// InternalAddrType is the address type for all keys within branch 1
	// (change addresses).
	InternalAddrType AddressType
}

var (
	// KeyScopeBIP0049Plus is the key scope of our modified BIP0049
	// derivation. We say this is BIP0049 "plus", as we'll actually use
	// p2wkh change all change addresses.
	KeyScopeBIP0049Plus = KeyScope{
		Purpose: 49,
		Coin:    0,
	}

	// KeyScopeBIP0084 is the key scope for BIP0084 derivation. BIP0084
	// will be used to derive all p2wkh addresses.
	KeyScopeBIP0084 = KeyScope{
		Purpose: 84,
		Coin:    0,
	}

	// KeyScopeBIP0044 is the key scope for BIP0044 derivation. Legacy
	// wallets will only be able to use this key scope, and no keys beyond
	// it.
	KeyScopeBIP0044 = KeyScope{
		Purpose: 44,
		Coin:    0,
	}

	// ScopeAddrMap is a map from the default key scopes to the scope
	// address schema for each scope type. This will be consulted during
	// the initial creation of the root key manager.
	ScopeAddrMap = map[KeyScope]ScopeAddrSchema{
		KeyScopeBIP0049Plus: {
			ExternalAddrType: NestedWitnessPubKey,
			InternalAddrType: WitnessPubKey,
		},
		KeyScopeBIP0084: {
			ExternalAddrType: WitnessPubKey,
			InternalAddrType: WitnessPubKey,
		},
		KeyScopeBIP0044: {
			InternalAddrType: PubKeyHash,
			ExternalAddrType: PubKeyHash,
		},
	}
)
