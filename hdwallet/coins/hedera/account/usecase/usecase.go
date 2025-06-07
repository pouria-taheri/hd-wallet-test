package usecase

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	domain "git.mazdax.tech/blockchain/hdwallet/coins/hedera/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip32"
	"fmt"
)

type useCase struct {
	Config            domain.Config
	coin              string
	logger            logger.Logger
	client            *hedera.Client
	operatorAccountID hedera.AccountID
	operatorPrivKey   hedera.PrivateKey
}

func New(logger logger.Logger, registry configcore.Registry,
	secureConfig config.SecureConfig, coin string) domain.UseCase {
	uc := &useCase{
		logger: logger,
		coin:   coin,
	}
	if err := registry.Unmarshal(&uc.Config); err != nil {
		panic(err)
	}
	uc.Config.Initialize(secureConfig)

	// Hedera-specific initialization
	// Read network and operator credentials from config
	network := uc.Config.ChainType // e.g., "mainnet" or "testnet"
	operatorIDStr := uc.Config.SecureConfig.OperatorAccountID // Add this field to your config
	operatorKeyStr := uc.Config.SecureConfig.OperatorPrivateKey // Add this field to your config

	fmt.Println("OperatorAccountID:", operatorIDStr)
	fmt.Println("OperatorPrivateKey:", operatorKeyStr)
	fmt.Println("ChainType:", network)

	var client *hedera.Client
	if network == "mainnet" {
		client = hedera.ClientForMainnet()
	} else {
		client = hedera.ClientForTestnet()
	}
	operatorAccountID, err := hedera.AccountIDFromString(operatorIDStr)
	if err != nil {
		panic(err)
	}
	operatorPrivKey, err := hedera.PrivateKeyFromString(operatorKeyStr)
	if err != nil {
		panic(err)
	}
	client.SetOperator(operatorAccountID, operatorPrivKey)

	uc.client = client
	uc.operatorAccountID = operatorAccountID
	uc.operatorPrivKey = operatorPrivKey

	return uc
}

// DeriveHederaKey derives a Hedera Ed25519 private key from a mnemonic and BIP44 path.
func DeriveHederaKey(mnemonic string, account, change, addressIndex uint32) (hedera.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	// BIP44 path: m/44'/3030'/account'/change/address_index
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 3030)
	accountKey, _ := coinType.NewChildKey(bip32.FirstHardenedChild + account)
	changeKey, _ := accountKey.NewChildKey(change)
	addressKey, _ := changeKey.NewChildKey(addressIndex)
	// Ed25519 private key from final key
	privKey, err := hedera.PrivateKeyFromBytes(addressKey.Key)
	if err != nil {
		return hedera.PrivateKey{}, err
	}
	return privKey, nil
}

// DeriveHederaPublicKey derives the public key from mnemonic and BIP44 path.
func DeriveHederaPublicKey(mnemonic string, account, change, addressIndex uint32) (string, error) {
	privKey, err := DeriveHederaKey(mnemonic, account, change, addressIndex)
	if err != nil {
		return "", err
	}
	pubKey := privKey.PublicKey()
	return pubKey.String(), nil
}

// CreateAccountOnChain creates a Hedera account using the provided public key and returns the new account ID.
func CreateAccountOnChain(pubKey string, initialBalance int64) (string, error) {
	// Operator credentials (should be securely stored in production)
	operatorAccountID, err := hedera.AccountIDFromString("0.0.9201078")
	if err != nil {
		return "", err
	}
	operatorPrivKey, err := hedera.PrivateKeyFromString("1b12d89d7e8f714b29034dc6e976b5ce2d50174d5716ccb32e5498344a9892fb")
	if err != nil {
		return "", err
	}

	// Use mainnet client for production
	client := hedera.ClientForMainnet()
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

func (u *useCase) Coin() string {
	return "hedera"
}

func (u *useCase) GetAccount(request ad.Request) (*ad.Account, error) {
	// Derive public key from mnemonic and path
	pubKeyStr, err := DeriveHederaPublicKey(u.Config.Mnemonic, request.Account, request.Branch, request.Index)
	if err != nil {
		return nil, err
	}
	// For now, use the public key as the address (Hedera account ID is created on-chain)
	account := &ad.Account{
		Index:     request.Index,
		Type:      0, // custom type for Hedera
		Address:   pubKeyStr,
		PublicKey: []byte(pubKeyStr),
		Ed25519PublicKey: []byte(pubKeyStr),
	}
	return account, nil
} 