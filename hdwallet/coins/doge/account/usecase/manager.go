package usecase

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"sync"
)

const (
	// MaxAccountNum is the maximum allowed account number.  This value was
	// chosen because accounts are hardened children and therefore must not
	// exceed the hardened child range of extended keys and it provides a
	// reserved account at the top of the range for supporting imported
	// addresses.
	MaxAccountNum = hdkeychain.HardenedKeyStart - 2 // 2^31 - 2

	// MaxAddressesPerAccount is the maximum allowed number of addresses
	// per account number.  This value is based on the limitation of the
	// underlying hierarchical deterministic key derivation.
	MaxAddressesPerAccount = hdkeychain.HardenedKeyStart - 1

	// ImportedAddrAccount is the account number to use for all imported
	// addresses.  This is useful since normal accounts are derived from
	// the root hierarchical deterministic key and imported addresses do
	// not fit into that model.
	ImportedAddrAccount = MaxAccountNum + 1 // 2^31 - 1

	// ImportedAddrAccountName is the name of the imported account.
	ImportedAddrAccountName = "imported"

	// DefaultAccountNum is the number of the default account.
	DefaultAccountNum = 0

	// defaultAccountName is the initial name of the default account.  Note
	// that the default account may be renamed and is not a reserved name,
	// so the default account might not be named "default" and non-default
	// accounts may be named "default".
	//
	// Account numbers never change, so the DefaultAccountNum should be
	// used to refer to (and only to) the default account.
	defaultAccountName = "default"

	// The hierarchy described by BIP0043 is:
	//  m/<purpose>'/*
	// This is further extended by BIP0044 to:
	//  m/44'/<coin type>'/<account>'/<branch>/<address index>
	//
	// The branch is 0 for external addresses and 1 for internal addresses.

	// maxCoinType is the maximum allowed coin type used when structuring
	// the BIP0044 multi-account hierarchy.  This value is based on the
	// limitation of the underlying hierarchical deterministic key
	// derivation.
	maxCoinType = hdkeychain.HardenedKeyStart - 1

	// ExternalBranch is the child number to use when performing BIP0044
	// style hierarchical deterministic key derivation for the external
	// branch.
	ExternalBranch uint32 = 0

	// InternalBranch is the child number to use when performing BIP0044
	// style hierarchical deterministic key derivation for the internal
	// branch.
	InternalBranch uint32 = 1

	// saltSize is the number of bytes of the salt used when hashing
	// private passphrases.
	saltSize = 32
)

// addrKey is used to uniquely identify an address even when those addresses
// would end up being the same bitcoin address (as is the case for
// pay-to-pubkey and pay-to-pubkey-hash style of addresses).
type addrKey string

// accountInfo houses the current state of the internal and external branches
// of an account along with the extended keys needed to derive new keys.  It
// also handles locking by keeping an encrypted version of the serialized
// private extended key so the unencrypted versions can be cleared from memory
// when the address manager is locked.
type accountInfo struct {
	// The account key is used to derive the branches which in turn derive
	// the internal and external addresses.  The accountKeyPriv will be nil
	// when the address manager is locked.
	acctKeyPriv *hdkeychain.ExtendedKey
	acctKeyPub  *hdkeychain.ExtendedKey

	// The external branch is used for all addresses which are intended for
	// external use.
	nextExternalIndex uint32
	lastExternalAddr  ad.ManagedAddress

	// The internal branch is used for all adddresses which are only
	// intended for internal wallet use such as change addresses.
	nextInternalIndex uint32
	lastInternalAddr  ad.ManagedAddress
}

// AccountProperties contains properties associated with each account, such as
// the account name, number, and the nubmer of derived and imported keys.
type AccountProperties struct {
	AccountNumber    uint32
	AccountName      string
	ExternalKeyCount uint32
	InternalKeyCount uint32
	ImportedKeyCount uint32
}

// Manager represents a concurrency safe crypto currency address manager and
// key store.
type Manager struct {
	mtx sync.RWMutex

	// scopedManager is a mapping of scope of scoped manager, the manager
	// itself loaded into memory.
	scopedManagers map[ad.KeyScope]*ScopedKeyManager

	externalAddrSchemas map[ad.AddressType][]ad.KeyScope
	internalAddrSchemas map[ad.AddressType][]ad.KeyScope

	chainParams *chaincfg.Params
}

// newManager returns a new locked address manager with the given parameters.
func newManager(chainParams *chaincfg.Params, scopedManagers map[ad.KeyScope]*ScopedKeyManager) *Manager {

	m := &Manager{
		chainParams:         chainParams,
		scopedManagers:      scopedManagers,
		externalAddrSchemas: make(map[ad.AddressType][]ad.KeyScope),
		internalAddrSchemas: make(map[ad.AddressType][]ad.KeyScope),
	}

	for _, sMgr := range m.scopedManagers {
		externalType := sMgr.AddrSchema().ExternalAddrType
		internalType := sMgr.AddrSchema().InternalAddrType
		scope := sMgr.Scope()

		m.externalAddrSchemas[externalType] = append(
			m.externalAddrSchemas[externalType], scope,
		)
		m.internalAddrSchemas[internalType] = append(
			m.internalAddrSchemas[internalType], scope,
		)
	}

	return m
}

// checkBranchKeys ensures deriving the extended keys for the internal and
// external branches given an account key does not result in an invalid child
// error which means the chosen seed is not usable.  This conforms to the
// hierarchy described by the BIP0044 family so long as the account key is
// already derived accordingly.
//
// In particular this is the hierarchical deterministic extended key path:
//   m/purpose'/<coin type>'/<account>'/<branch>
//
// The branch is 0 for external addresses and 1 for internal addresses.
func checkBranchKeys(acctKey *hdkeychain.ExtendedKey) error {
	// Derive the external branch as the first child of the account key.
	if _, err := acctKey.Child(ExternalBranch); err != nil {
		return err
	}

	// Derive the external branch as the second child of the account key.
	_, err := acctKey.Child(InternalBranch)
	return err
}

// LoadManager returns a new address manager that results from loading it from
// the passed opened database.  The public passphrase is required to decrypt
// the public keys.
func LoadManager(masterNode *hdkeychain.ExtendedKey, chainParams *chaincfg.Params, scopes []ad.KeyScope,
	ScopeAddrMap map[ad.KeyScope]ad.ScopeAddrSchema, log logger.Logger) (*Manager, error) {

	// Next, we'll need to load all known manager scopes from disk. Each
	// scope is on a distinct top-level path within our HD key chain.
	scopedManagers := make(map[ad.KeyScope]*ScopedKeyManager)

	for _, scope := range scopes {
		scopedManagers[scope] = &ScopedKeyManager{
			scope:      scope,
			addrSchema: ScopeAddrMap[scope],
			addrs:      make(map[addrKey]ad.ManagedAddress),
			acctInfo:   make(map[uint32]*accountInfo),
			logger:     log,
		}
	}

	// Create new address manager with the given parameters.  Also,
	// override the defaults for the additional fields which are not
	// specified in the call to new with the values loaded from the
	// database.
	mgr := newManager(chainParams, scopedManagers)

	for _, scopedManager := range scopedManagers {
		// Derive the cointype key according to the passed scope.
		coinTypeKeyPriv, err := deriveCoinTypeKey(masterNode,
			scopedManager.Scope())
		if err != nil {
			log.With(logger.Field{
				"submodule": "manager usecase",
				"section":   "loadManager",
				"error":     err,
				"scope":     scopedManager.Scope(),
			}).ErrorF("cannot get coin type of extended key")
			return nil, errors.New().
				WithMessage("failed to derive coin type extended key")
		}
		coinTypeKeyPub, err := coinTypeKeyPriv.Neuter()
		if err != nil {
			log.With(logger.Field{
				"submodule": "manager usecase",
				"section":   "loadManager",
				"error":     err,
			}).ErrorF("cannot convert coin type private key to public")
			return nil, errors.New().
				WithMessage("failed to convert coin type private key")
		}

		scopedManager.coinTypeKey = coinTypeKeyPriv
		scopedManager.coinTypeKeyPub = coinTypeKeyPub

		scopedManager.rootManager = mgr
	}
	mgr.scopedManagers = scopedManagers

	return mgr, nil
}

// deriveCoinTypeKey derives the cointype key which can be used to derive the
// extended key for an account according to the hierarchy described by BIP0044
// given the coin type key.
//
// In particular this is the hierarchical deterministic extended key path:
// m/purpose'/<coin type>'
func deriveCoinTypeKey(masterNode *hdkeychain.ExtendedKey, scope ad.KeyScope) (*hdkeychain.ExtendedKey, error) {

	// Enforce maximum coin type.
	if scope.Coin > maxCoinType {
		return nil, errors.New().
			WithMessage("Coin type is too high")
	}

	// The hierarchy described by BIP0043 is:
	//  m/<purpose>'/*
	//
	// This is further extended by BIP0044 to:
	//  m/44'/<coin type>'/<account>'/<branch>/<address index>
	//
	// However, as this is a generic key store for any family for BIP0044
	// standards, we'll use the custom scope to govern our key derivation.
	//
	// The branch is 0 for external addresses and 1 for internal addresses.

	// Derive the purpose key as a child of the master node.
	purpose, err := masterNode.Child(scope.Purpose + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Derive the coin type key as a child of the purpose key.
	coinTypeKey, err := purpose.Child(scope.Coin + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	return coinTypeKey, nil
}

// deriveAccountKey derives the extended key for an account according to the
// hierarchy described by BIP0044 given the master node.
//
// In particular this is the hierarchical deterministic extended key path:
//   m/purpose'/<coin type>'/<account>'
func (m *Manager) deriveAccountKey(coinTypeKey *hdkeychain.ExtendedKey,
	account uint32) (*hdkeychain.ExtendedKey, error) {

	// Enforce maximum account number.
	if account > MaxAccountNum {
		return nil, errors.New().
			WithMessage("account number is too high")
	}

	// Derive the account key as a child of the coin type key.
	return coinTypeKey.Child(account + hdkeychain.HardenedKeyStart)
}

// GetAccount creates a new key scoped for a target manager's scope.
// This partitions key derivation for a particular purpose+coin tuple, allowing
// multiple address derivation schems to be maintained concurrently.
func (m *Manager) GetAccount(scope ad.KeyScope, request ad.Request) (*ad.Account, error) {

	scopeMgr, ok := m.scopedManagers[scope]
	if !ok {
		return nil, errors.New().
			WithMessage("no manager found for scope")
	}

	accKey, err := m.deriveAccountKey(scopeMgr.coinTypeKey, request.Account)
	if err != nil {
		scopeMgr.logger.With(logger.Field{
			"submodule":    "manager usecase",
			"section":      "get account",
			"error":        err,
			"request body": request,
		}).ErrorF("cannot get account extended key from coinType and account")
		return nil, err
	}

	// Choose the account key to used based on whether the
	// request is for private.
	if !request.Private {
		accKey, err = accKey.Neuter()
		if err != nil {
			scopeMgr.logger.With(logger.Field{
				"submodule":    "manager usecase",
				"section":      "get account",
				"error":        err,
				"request body": request,
			}).ErrorF("cannot neuter of account extended key")
			return nil, err
		}
	}

	managedAddr, err := scopeMgr.keyToManaged(accKey, request.Account,
		request.Branch, request.Index)
	if err != nil {
		scopeMgr.logger.With(logger.Field{
			"submodule":    "manager usecase",
			"section":      "get account",
			"error":        err,
			"request body": request,
		}).ErrorF("cannot get manage address from scope")
		return nil, err
	}

	acc := &ad.Account{
		Id:      uint64(request.Account),
		Index:   request.Index,
		Type:    managedAddr.AddrType().Uint8(),
		Address: managedAddr.Address().String(),
	}

	return acc, nil
}

// GetManagedAddress creates a new key scoped for a target manager's scope.
// This partitions key derivation for a particular purpose+coin tuple, allowing
// multiple address derivation schems to be maintained concurrently.
func (m *Manager) GetManagedAddress(scope ad.KeyScope, request ad.Request) (ad.ManagedAddress, error) {

	scopeMgr, ok := m.scopedManagers[scope]
	if !ok {
		return nil, errors.New().
			WithMessage("no manager found for scope")
	}

	accKey, err := m.deriveAccountKey(scopeMgr.coinTypeKey, request.Account)
	if err != nil {
		scopeMgr.logger.With(logger.Field{
			"submodule":    "manager usecase",
			"section":      "get managed address",
			"error":        err,
			"request body": request,
		}).ErrorF("cannot get account extended key from coinType and account")
		return nil, err
	}

	// Choose the account key to used based on whether the
	// request is for private.
	if !request.Private {
		accKey, err = accKey.Neuter()
		if err != nil {
			scopeMgr.logger.With(logger.Field{
				"submodule":    "manager usecase",
				"section":      "get managed address",
				"error":        err,
				"request body": request,
			}).ErrorF("cannot neuter of account extended key")
			return nil, err
		}
	}

	managedAddr, err := scopeMgr.keyToManaged(accKey, request.Account,
		request.Branch, request.Index)
	if err != nil {
		scopeMgr.logger.With(logger.Field{
			"submodule":    "manager usecase",
			"section":      "get managed address",
			"error":        err,
			"request body": request,
		}).ErrorF("cannot get manage address from scope")
		return nil, err
	}

	return managedAddr, nil
}
