package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

type SecureConfigDetail struct {
	Coin     string
	FilePath string
	Data     []byte
}

type SecureConfig struct {
	Mnemonic string
	Seed     string
	Password string

	detail SecureConfigDetail

	// Hedera-specific fields
	OperatorAccountID  string
	OperatorPrivateKey string
	ChainType          string
}

func (c *SecureConfig) SetDetail(detail SecureConfigDetail) {
	c.detail = detail
}

func (c *SecureConfig) GetDetail() SecureConfigDetail {
	return c.detail
}

func (c *SecureConfig) getSeedFromSeedString() ([]byte, error) {
	c.Seed = strings.TrimSpace(strings.ToLower(c.Seed))

	seed, err := hex.DecodeString(c.Seed)
	if err != nil || len(seed) < hdkeychain.MinSeedBytes ||
		len(seed) > hdkeychain.MaxSeedBytes {

		return nil, errors.New(fmt.Sprintf("Invalid seed specified. "+
			"Must be a hexadecimal value that is at least %d bits and at most "+
			"%d bits", hdkeychain.MinSeedBytes*8, hdkeychain.MaxSeedBytes*8))
	}
	return seed, nil
}

func (c *SecureConfig) GetSeed() ([]byte, error) {
	if c.Mnemonic != "" {
		seed := bip39.NewSeed(c.Mnemonic, "")
		return seed, nil
	}
	return c.getSeedFromSeedString()
}

func (c *SecureConfig) GetSeedWithErrorChecking() ([]byte, error) {
	if c.Mnemonic != "" {
		seed, err := bip39.NewSeedWithErrorChecking(c.Mnemonic, c.Password)
		return seed, err
	}
	return c.getSeedFromSeedString()
}

func (c *SecureConfig) SetRandomDta(bitSize int) {
	ent, err := bip39.NewEntropy(bitSize)
	if err != nil {
		panic(err)
	}
	mnemonic, err := bip39.NewMnemonic(ent)
	if err != nil {
		panic(err)
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}
	c.Mnemonic = mnemonic
	c.Seed = hex.EncodeToString(seed)
}
