package wallet

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/hedera/account/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type HederaWallet struct {
	Config    *domain.Config
	Logger    logger.Logger
	MasterKey *bip32.Key
}

type HederaWalletInterface interface {
	DerivePrivateKey(account, change, addressIndex uint32) (hedera.PrivateKey, error)
	DerivePublicKey(account, change, addressIndex uint32) (string, error)
}

func NewHederaHdWallet(config *domain.Config, logger logger.Logger) *HederaWallet {
	seed := bip39.NewSeed(config.Mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		logger.Error("Failed to create master key", err)
		return nil
	}
	return &HederaWallet{
		Config:    config,
		Logger:    logger,
		MasterKey: masterKey,
	}
}

// DerivePrivateKey derives a Hedera Ed25519 private key from the wallet's master key using BIP44 path.
func (w *HederaWallet) DerivePrivateKey(account, change, addressIndex uint32) (hedera.PrivateKey, error) {
	purpose, err := w.MasterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 3030)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	accountKey, err := coinType.NewChildKey(bip32.FirstHardenedChild + account)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	changeKey, err := accountKey.NewChildKey(change)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	addressKey, err := changeKey.NewChildKey(addressIndex)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	privKey, err := hedera.PrivateKeyFromBytes(addressKey.Key)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	return privKey, nil
}

// DerivePublicKey derives the public key from the wallet's master key using BIP44 path.
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
	pubKeyStr, err := w.DerivePublicKey(account, change, addressIndex)
	if err != nil {
		return nil, err
	}
	return &ad.Account{
		Index:            addressIndex,
		Type:             0, // custom type for Hedera
		Address:          pubKeyStr,
		PublicKey:        []byte(pubKeyStr),
		Ed25519PublicKey: []byte(pubKeyStr),
	}, nil
}
