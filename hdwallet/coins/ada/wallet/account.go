package wallet

import (
	"bytes"
	"encoding/gob"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/crypto"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/db"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
)

type cardanoWallet struct {
	domain.MasterKey
	db.MasterWallet
	domain.WalletMainAcc
}

func (w *cardanoWallet) Coin() string {
	return "ADA"
}

func (w *cardanoWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	addr := w.Wallet.GetAddress(request.Account, request.Index)
	w.DB.SaveWallet(w.Wallet)

	return &ad.Account{
		Address: addr.String(),
	}, nil
}

func (w *cardanoWallet) GetAddresses() ([]byte, error) {
	addresses := w.Wallet.Addresses()

	var payload bytes.Buffer
	err := gob.NewEncoder(&payload).Encode(addresses)
	if err != nil {
		return nil, err
	}
	return payload.Bytes(), nil
}

func (w *cardanoWallet) GetWallet() ([]byte, error) {
	var payload bytes.Buffer
	err := gob.NewEncoder(&payload).Encode(w.Wallet)
	if err != nil {
		return nil, err
	}
	return payload.Bytes(), nil
}

func (w *cardanoWallet) NewAddress() ([]byte, error) {
	address := w.Wallet.AddAddress()
	w.DB.SaveWallet(w.Wallet)

	var payload bytes.Buffer
	err := gob.NewEncoder(&payload).Encode(address)
	if err != nil {
		return nil, err
	}
	return payload.Bytes(), nil
}

func (w *cardanoWallet) AccPrv(accountId, index uint32) (crypto.ExtendedSigningKey, error) {
	return nil, nil
}
