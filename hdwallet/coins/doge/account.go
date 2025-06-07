package doge

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/doge/account"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(registry configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) ad.UseCase {
	uc := account.NewUseCase(registry, secureConfig, "doge", logger)

	var chainParams = uc.GetChainParams()
	var (
		// keyScopeBIP0049Plus is the key scope of our modified BIP0049
		// derivation. We say this is BIP0049 "plus", as we'll actually use
		// p2wkh change all change addresses.
		keyScopeBIP0049Plus = ad.KeyScope{
			Purpose: 49,
			Coin:    chainParams.HDCoinType,
		}
		// keyScopeBIP0044 is the key scope for BIP0044 derivation. Legacy
		// wallets will only be able to use this key scope, and no keys beyond
		// it.
		keyScopeBIP0044 = ad.KeyScope{
			Purpose: 44,
			Coin:    chainParams.HDCoinType,
		}
		// defaultKeyScopes is the set of default key scopes that will be
		// created by the root manager upon initial creation.
		defaultKeyScopes = []ad.KeyScope{
			keyScopeBIP0049Plus,
			keyScopeBIP0044,
		}
		// scopeAddrMap is a map from the default key scopes to the scope
		// address schema for each scope type. This will be consulted during
		// the initial creation of the root key manager.
		scopeAddrMap = map[ad.KeyScope]ad.ScopeAddrSchema{
			keyScopeBIP0049Plus: {
				ExternalAddrType: ad.NestedWitnessPubKey,
				InternalAddrType: ad.WitnessPubKey,
			},
			keyScopeBIP0044: {
				InternalAddrType: ad.PubKeyHash,
				ExternalAddrType: ad.PubKeyHash,
			},
		}
	)
	uc.SetChainParams(chainParams)

	uc.Initialize(defaultKeyScopes, scopeAddrMap)
	return uc
}
