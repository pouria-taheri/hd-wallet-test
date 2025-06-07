package db

import (
	"encoding/json"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/address"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/crypto"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tyler-smith/go-bip39"
	"reflect"
)

const (
	EntropySizeInBits         = 160
	purposeIndex       uint32 = 1852 + 0x80000000
	coinTypeIndex      uint32 = 1815 + 0x80000000
	accountIndex       uint32 = 0x80000000
	externalChainIndex uint32 = 0x0
	walleIDAlphabet           = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type Wallet struct {
	ID      string
	Name    string
	Skeys   []crypto.ExtendedSigningKey
	pkeys   []crypto.ExtendedVerificationKey
	rootKey crypto.ExtendedSigningKey
	coinKey crypto.ExtendedSigningKey
	Network address.Network
}

func (w *Wallet) SetNetwork(net address.Network) {
	w.Network = net
}

// AddAddress generates a new payment Address and adds it to the wallet.
func (w *Wallet) AddAddress() address.Address {
	index := uint32(len(w.Skeys))
	accountKey := crypto.DeriveSigningKey(w.rootKey, accountIndex)
	chainKey := crypto.DeriveSigningKey(accountKey, externalChainIndex)
	addrKey := crypto.DeriveSigningKey(chainKey, index)

	if !w.addressExists(addrKey) {
		w.Skeys = append(w.Skeys, addrKey)
	}
	return address.NewEnterpriseAddress(addrKey.ExtendedVerificationKey(), w.Network)
}

// GetAddress get generated address for specific accountId and index.
func (w *Wallet) GetAddress(accountId, index uint32) address.Address {
	if accountId == 0 && index == 0 {
		return address.NewEnterpriseAddress(w.Skeys[0].ExtendedVerificationKey(), w.Network)
	}
	accountKey := crypto.DeriveSigningKey(w.rootKey, accountIndex+accountId)
	chainKey := crypto.DeriveSigningKey(accountKey, externalChainIndex)
	addrKey := crypto.DeriveSigningKey(chainKey, index)

	if !w.addressExists(addrKey) {
		w.Skeys = append(w.Skeys, addrKey)
	}
	return address.NewEnterpriseAddress(addrKey.ExtendedVerificationKey(), w.Network)
}

func (w *Wallet) addressExists(addr crypto.ExtendedSigningKey) bool {
	for _, key := range w.Skeys {
		if reflect.DeepEqual(key.ExtendedVerificationKey(), addr.ExtendedVerificationKey()) {
			return true
		}
	}
	return false
}

// AddAddresses returns all wallet's addresss.
func (w *Wallet) Addresses() []address.Address {
	addresses := make([]address.Address, len(w.Skeys))
	for _, key := range w.Skeys {
		addresses = append(addresses, address.NewEnterpriseAddress(key.ExtendedVerificationKey(), w.Network))
	}
	return addresses
}

func newWalletID() string {
	id, _ := gonanoid.Generate(walleIDAlphabet, 10)
	return "wallet_" + id
}

func NewWallet(name, password string, entropy []byte) *Wallet {
	wallet := &Wallet{Name: name, ID: newWalletID()}
	wallet.Skeys = []crypto.ExtendedSigningKey{}
	rootKey := crypto.NewExtendedSigningKey(entropy, password)
	purposeKey := crypto.DeriveSigningKey(rootKey, purposeIndex)
	coinKey := crypto.DeriveSigningKey(purposeKey, coinTypeIndex)
	accountKey := crypto.DeriveSigningKey(coinKey, accountIndex)
	chainKey := crypto.DeriveSigningKey(accountKey, externalChainIndex)
	addr0Key := crypto.DeriveSigningKey(chainKey, 0)
	wallet.rootKey = chainKey
	wallet.coinKey = coinKey
	wallet.Skeys = []crypto.ExtendedSigningKey{addr0Key}
	return wallet
}

type walletDump struct {
	ID      string
	Name    string
	Keys    []crypto.ExtendedSigningKey
	RootKey crypto.ExtendedSigningKey
}

func (w *Wallet) Marshal() ([]byte, error) {
	wd := &walletDump{
		ID:      w.ID,
		Name:    w.Name,
		Keys:    w.Skeys,
		RootKey: w.rootKey,
	}
	bytes, err := json.Marshal(wd)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (w *Wallet) Unmarshal(bytes []byte) error {
	wd := &walletDump{}
	err := json.Unmarshal(bytes, wd)
	if err != nil {
		return err
	}
	w.ID = wd.ID
	w.Name = wd.Name
	w.Skeys = wd.Keys
	w.rootKey = wd.RootKey
	return nil
}

var NewEntropy = func(bitSize int) []byte {
	entropy, _ := bip39.NewEntropy(bitSize)
	return entropy
}
