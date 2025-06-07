package domain

import "fmt"

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

// ScopedIndex is a tuple of KeyScope and child Index. This is used to compactly
// identify a particular child key, when the account and branch can be inferred
// from context.
type ScopedIndex struct {
	// Scope is the BIP44 account' used to derive the child key.
	Scope KeyScope

	// Index is the BIP44 address_index used to derive the child key.
	Index uint32
}

// String returns a human readable version describing the keypath encapsulated
// by the target key scope.
func (k *KeyScope) String() string {
	return fmt.Sprintf("m/%v'/%v'", k.Purpose, k.Coin)
}
