package wallet

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

type terraWallet struct {
	domain.MasterKey
	domain.WalletMainAcc
}

func (w *terraWallet) Coin() string {
	return "LUNA"
}

func (w *terraWallet) GetAccount(request ad.Request) (*ad.Account, error) {
	var memoBaseValue uint64 = 78563452053
	path := hd.NewFundraiserParams(request.Account, 330, 0)

	prv, errPrv := hd.DerivePrivateKeyForPath(w.Secret, w.ChainCode, path.String())
	if errPrv != nil {
		return nil, errPrv
	}

	masterAddress, err := ConvertAndEncode("terra", tmsecp256k1.GenPrivKeySecp256k1(w.PrivateKey).PubKey().Address())
	if err != nil {
		return nil, err
	}

	address, err := ConvertAndEncode("terra", tmsecp256k1.GenPrivKeySecp256k1(prv).PubKey().Address())
	if err != nil {
		return nil, err
	}

	return &ad.Account{
		Master:  masterAddress,
		Address: address,
		Memo:    fmt.Sprintf("%v", memoBaseValue+uint64(request.Account)),
	}, nil
}

func (w *terraWallet) AccPrv(path *hd.BIP44Params) (tmsecp256k1.PrivKey, error) {
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

//ConvertAndEncode converts from a base64 encoded byte string to base32 encoded byte string and then to bech32
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed: %v", err)
	}
	return bech32.Encode(hrp, converted)

}

//DecodeAndConvert decodes a bech32 encoded string and converts to base64 encoded bytes
func DecodeAndConvert(bech string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(bech)
	if err != nil {
		return "", nil, fmt.Errorf("decoding bech32 failed: %v", err)
	}
	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, fmt.Errorf("decoding bech32 failed: %v", err)
	}
	return hrp, converted, nil
}
