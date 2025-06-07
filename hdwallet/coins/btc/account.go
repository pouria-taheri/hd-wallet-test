package btc

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/account"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

var (
	// KeyScopeBIP0049Plus is the key scope of our modified BIP0049
	// derivation. We say this is BIP0049 "plus", as we'll actually use
	// p2wkh change all change addresses.
	KeyScopeBIP0049Plus = ad.KeyScope{
		Purpose: 49,
		Coin:    0,
	}

	// KeyScopeBIP0084 is the key scope for BIP0084 derivation. BIP0084
	// will be used to derive all p2wkh addresses.
	KeyScopeBIP0084 = ad.KeyScope{
		Purpose: 84,
		Coin:    0,
	}

	// KeyScopeBIP0044 is the key scope for BIP0044 derivation. Legacy
	// wallets will only be able to use this key scope, and no keys beyond
	// it.
	KeyScopeBIP0044 = ad.KeyScope{
		Purpose: 44,
		Coin:    0,
	}

	// DefaultKeyScopes is the set of default key scopes that will be
	// created by the root manager upon initial creation.
	DefaultKeyScopes = []ad.KeyScope{
		KeyScopeBIP0049Plus,
		KeyScopeBIP0084,
		KeyScopeBIP0044,
	}

	// ScopeAddrMap is a map from the default key scopes to the scope
	// address schema for each scope type. This will be consulted during
	// the initial creation of the root key manager.
	ScopeAddrMap = map[ad.KeyScope]ad.ScopeAddrSchema{
		KeyScopeBIP0049Plus: {
			ExternalAddrType: ad.NestedWitnessPubKey,
			InternalAddrType: ad.WitnessPubKey,
		},
		KeyScopeBIP0084: {
			ExternalAddrType: ad.WitnessPubKey,
			InternalAddrType: ad.WitnessPubKey,
		},
		KeyScopeBIP0044: {
			InternalAddrType: ad.PubKeyHash,
			ExternalAddrType: ad.PubKeyHash,
		},
	}
)

func NewAccountUseCase(registry configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) ad.UseCase {

	uc := account.NewUseCase(registry, secureConfig, "btc", logger)
	uc.Initialize(DefaultKeyScopes, ScopeAddrMap)
	return uc
}
