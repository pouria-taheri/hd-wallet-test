package wallet

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/atom/domain"
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

type cosmosWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *cosmosWallet) Coin() string {
	return "ATOM"
}

func (w *cosmosWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 35498922654
	path := hd.NewFundraiserParams(request.Account, 118, 0)

	prv, errPrv := hd.DerivePrivateKeyForPath(w.Secret, w.ChainCode, path.String())
	if errPrv != nil {
		return nil, errPrv
	}

	privateKey := secp256k1.GenPrivKeyFromSecret(prv)
	addr := sdk.AccAddress(privateKey.PubKey().Address())

	pk, _ := btcec.ParsePubKey(privateKey.PubKey().Bytes(), btcec.S256())

	return &ad.Account{
		Master:    sdk.AccAddress(secp256k1.GenPrivKeyFromSecret(w.PrivateKey).PubKey().Address()).String(),
		Address:   addr.String(),
		PublicKey: pk.SerializeCompressed(),
		Memo:      fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
	}, nil
}

func (w *cosmosWallet) AccPrv(path *hd.BIP44Params) (tmsecp256k1.PrivKey, error) {
	prv, errPrv := hd.DerivePrivateKeyForPath(w.Secret, w.ChainCode, path.String())
	if errPrv != nil {
		return tmsecp256k1.PrivKey{}, errPrv
	}

	privateKey := tmsecp256k1.GenPrivKeySecp256k1(prv)
	return privateKey, nil
}

func convertTo32Byte(data []byte) [32]byte {
	derivedKey := make([]byte, 32)
	copy(derivedKey, data[:])
	x2 := [32]byte{}
	copy(x2[32-len(data):], data)

	return x2
}
