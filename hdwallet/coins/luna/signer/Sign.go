package signer

import (
	"encoding/hex"
	"encoding/json"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

type Signer struct {
	wallet      domain.TerraWallet
	logger      logger.Logger
	NetworkType string
}

type Transaction struct {
	From             string `json:"from"`
	To               string `json:"to"`
	Value            int64  `json:"value"`
	Memo             string `json:"memo"`
	SignaturePayload string `json:"signature"`
}

type SignMsg struct {
	Transaction Transaction `json:"transaction"`
	Account     ad.Account  `json:"account"`
}

type StdSignature struct {
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
}

type signerConf struct {
	ChainType string
}

func (s Signer) Coin() string {
	return "luna"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {

	var signMsg SignMsg
	err := json.Unmarshal(request, &signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}

	tx := signMsg.Transaction

	path := hd.NewFundraiserParams(uint32(signMsg.Account.Id), 330, 0)
	prv, err := s.wallet.AccPrv(path)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     path,
		}).ErrorF("error in get account private key")
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	sk, err := prv.Sign([]byte(tx.SignaturePayload))
	if err != nil {
		return nil, err
	}

	signature := hex.EncodeToString(sk[:])
	publicKey := hex.EncodeToString(prv.PubKey().Bytes()[5:])

	res, err := json.Marshal(&StdSignature{
		Signature: signature,
		PublicKey: publicKey,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.TerraWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
		sc.ChainType,
	}
	return s
}
