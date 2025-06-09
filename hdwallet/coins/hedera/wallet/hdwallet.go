package wallet

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/hedera/account/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/ChainSafe/slip10"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/tyler-smith/go-bip39"
)

type HederaWallet struct {
	Config *domain.Config
	Logger logger.Logger
	Seed   []byte
}

type HederaWalletInterface interface {
	DerivePrivateKey(account, change, addressIndex uint32) (hedera.PrivateKey, error)
	DerivePublicKey(account, change, addressIndex uint32) (string, error)
}

func NewHederaHdWallet(config *domain.Config, logger logger.Logger) *HederaWallet {
	seed := bip39.NewSeed(config.Mnemonic, "")
	return &HederaWallet{
		Config: config,
		Logger: logger,
		Seed:   seed,
	}
}

// DerivePrivateKey derives a Hedera Ed25519 private key from the wallet's seed using SLIP-0010/BIP44 path.
func (w *HederaWallet) DerivePrivateKey(account, change, addressIndex uint32) (hedera.PrivateKey, error) {
	// BIP44 path: m/44'/3030'/account'/change/address_index
	path := []uint32{
		slip10.Hardened + 44,
		slip10.Hardened + 3030,
		slip10.Hardened + account,
		change,       // not hardened
		addressIndex, // not hardened
	}
	key, err := slip10.DeriveKeyFromPath(w.Seed, slip10.Ed25519(), path)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	privKey, err := hedera.PrivateKeyFromBytes(key.Key)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	return privKey, nil
}

// DerivePublicKey derives the public key from the wallet's seed using SLIP-0010/BIP44 path.
func (w *HederaWallet) DerivePublicKey(account, change, addressIndex uint32) (string, error) {
	privKey, err := w.DerivePrivateKey(account, change, addressIndex)
	if err != nil {
		return "", err
	}
	return privKey.PublicKey().String(), nil
}

// CreateAccountOnChain creates a Hedera account using the provided public key and returns the new account ID.
func (w *HederaWallet) CreateAccountOnChain(pubKey string, initialBalance int64) (string, error) {
	operatorAccountID, err := hedera.AccountIDFromString(w.Config.OperatorAccountID)
	if err != nil {
		return "", err
	}
	operatorPrivKey, err := hedera.PrivateKeyFromString(w.Config.OperatorPrivateKey)
	if err != nil {
		return "", err
	}
	var client *hedera.Client
	if w.Config.ChainType == "mainnet" {
		client = hedera.ClientForMainnet()
	} else {
		client = hedera.ClientForTestnet()
	}
	client.SetOperator(operatorAccountID, operatorPrivKey)
	userPubKey, err := hedera.PublicKeyFromString(pubKey)
	if err != nil {
		return "", err
	}
	tx, err := hedera.NewAccountCreateTransaction().
		SetKey(userPubKey).
		SetInitialBalance(hedera.HbarFrom(initialBalance, hedera.HbarUnits.Tinybar)).
		Execute(client)
	if err != nil {
		return "", err
	}
	receipt, err := tx.GetReceipt(client)
	if err != nil {
		return "", err
	}
	if receipt.AccountID == nil {
		return "", err
	}
	return receipt.AccountID.String(), nil
}

// GetAccount returns an ad.Account for the given derivation path.
func (w *HederaWallet) GetAccount(mnemonic string, account, change, addressIndex uint32) (*ad.Account, error) {
	privKey, err := w.DerivePrivateKey(account, change, addressIndex)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := privKey.PublicKey().Bytes()
	pubKeyStr := privKey.PublicKey().String()
	return &ad.Account{
		Index:            addressIndex,
		Type:             0, // custom type for Hedera
		Address:          pubKeyStr,
		PublicKey:        pubKeyBytes,
		Ed25519PublicKey: pubKeyBytes,
	}, nil
}
