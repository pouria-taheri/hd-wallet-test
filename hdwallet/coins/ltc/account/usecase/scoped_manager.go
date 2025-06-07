package usecase

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/ltcsuite/ltcutil/hdkeychain"
	"sync"
)

// ScopedKeyManager is a sub key manager under the main root key manager. The
// root key manager will handle the root HD key (m/), while each sub scoped key
// manager will handle the cointype key for a particular key scope
// (m/purpose'/cointype'). This abstraction allows higher-level applications
// built upon the root key manager to perform their own arbitrary key
// derivation, while still being protected under the encryption of the root key
// manager.
type ScopedKeyManager struct {
	// scope is the scope of this key manager. We can only generate keys
	// that are direct children of this scope.
	scope ad.KeyScope

	// addrSchema is the address schema for this sub manager. This will be
	// consulted when encoding addresses from derived keys.
	addrSchema ad.ScopeAddrSchema

	// coinTypeKey is key of master
	coinTypeKey *hdkeychain.ExtendedKey
	// coinTypeKeyPub is public key of master
	coinTypeKeyPub *hdkeychain.ExtendedKey

	// rootManager is a pointer to the root key manager. We'll maintain
	// this as we need access to the crypto encryption keys before we can
	// derive any new accounts of child keys of accounts.
	rootManager *Manager

	// addrs is a cached map of all the addresses that we currently
	// manager.
	addrs map[addrKey]ad.ManagedAddress

	// acctInfo houses information about accounts including what is needed
	// to generate deterministic chained keys for each created account.
	acctInfo map[uint32]*accountInfo

	mtx    sync.RWMutex
	logger logger.Logger
}

// Scope returns the exact KeyScope of this scoped key manager.
func (s *ScopedKeyManager) Scope() ad.KeyScope {
	return s.scope
}

// AddrSchema returns the set address schema for the target ScopedKeyManager.
func (s *ScopedKeyManager) AddrSchema() ad.ScopeAddrSchema {
	return s.addrSchema
}

// keyToManaged returns a new managed address for the provided derived key and
// its derivation path which consists of the account, branch, and index.
//
// The passed derivedKey is zeroed after the new address is created.
//
// This function MUST be called with the manager lock held for writes.
func (s *ScopedKeyManager) keyToManaged(derivedKey *hdkeychain.ExtendedKey,
	account, branch, index uint32) (ad.ManagedAddress, error) {

	var addrType ad.AddressType
	if branch == InternalBranch {
		addrType = s.addrSchema.InternalAddrType
	} else {
		addrType = s.addrSchema.ExternalAddrType
	}

	derivationPath := ad.DerivationPath{
		Account: account,
		Branch:  branch,
		Index:   index,
	}

	// Derive the appropriate branch key and ensure it is zeroed when done.
	branchKey, err := derivedKey.Child(branch)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "scope manager usecase",
			"error":     err,
			"branch":    branch,
		}).ErrorF("cannot get extended key")
		return nil, errors.New(err).WithMessage(fmt.Sprintf("failed to "+
			"derive extended key branch %d", branch))
	}
	defer branchKey.Zero() // Ensure branch key is zeroed when done.

	key, err := branchKey.Child(index)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "scope manager usecase",
			"error":     err,
			"index":     index,
		}).ErrorF("cannot generate child")
		return nil, errors.New(err).WithMessage(fmt.Sprintf("failed to "+
			"generate child %d", index))
	}
	key.SetNet(s.rootManager.chainParams)

	// Create a new managed address based on the public or private key
	// depending on whether the passed key is private.  Also, zero the key
	// after creating the managed address from it.
	ma, err := newManagedAddressFromExtKey(
		s, derivationPath, key, addrType,
	)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "scope manager usecase",
			"section":   "newManagedAddressFromExtKey",
			"error":     err,
		}).ErrorF("cannot get private key from extended key")
		return nil, err
	}

	if branch == InternalBranch {
		ma.internal = true
	}
	key.Zero()

	return ma, nil
}

// deriveKey returns either a public or private derived extended key based on
// the private flag for the given an account info, branch, and index.
func (s *ScopedKeyManager) deriveKey(acctInfo *accountInfo, branch,
	index uint32, private bool) (*hdkeychain.ExtendedKey, error) {

	// Choose the public or private extended key based on whether or not
	// the private flag was specified.  This, in turn, allows for public or
	// private child derivation.
	acctKey := acctInfo.acctKeyPub
	if private {
		acctKey = acctInfo.acctKeyPriv
	}

	// Derive and return the key.
	branchKey, err := acctKey.Child(branch)
	if err != nil {
		return nil, errors.New(err).
			WithMessage(fmt.Sprintf("failed to derive extended key branch"+
				" %d", branch))
	}

	addressKey, err := branchKey.Child(index)
	branchKey.Zero() // Zero branch key after it's used.
	if err != nil {
		return nil, errors.New(err).
			WithMessage(fmt.Sprintf("failed to derive child extended key "+
				"-- branch %d, child %d", branch, index))
	}

	return addressKey, nil
}
