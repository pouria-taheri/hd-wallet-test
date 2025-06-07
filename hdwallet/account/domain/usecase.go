package domain

type UseCase interface {
	Coin() string
	Initialize(scopes []KeyScope, scopeAddrMap map[KeyScope]ScopeAddrSchema)
	GetAccount(request Request) (*Account, error)
	GetManagedAddress(request Request) (ManagedAddress, error)

	GetChainType() string
}
