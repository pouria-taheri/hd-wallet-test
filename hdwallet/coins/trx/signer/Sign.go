package signer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/trx/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

type Signer struct {
	wallet domain.TronWallet
	logger logger.Logger
}

type SignMsg struct {
	Transaction *core.Transaction `json:"transaction"`
	Account     *ad.Account       `json:"account"`
}

type signerConf struct {
	NetworkId int64
}

func (s Signer) Coin() string {
	return "trx"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {

	signMsg := new(SignMsg)
	err := json.Unmarshal(request, signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}

	tx := signMsg.Transaction
	dPath := pathFromUserID(signMsg.Account.Id, signMsg.Account.Index)
	privateKey, err := s.wallet.AccPrv(dPath)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath,
		}).ErrorF("error in get account private key")
		return nil, err
	}

	tx.GetRawData().Timestamp = time.Now().UnixNano() / 1000000
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		s.logger.With(logger.Field{
			"error": err,
		}).ErrorF("error in get transaction raw data")
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	contractList := tx.GetRawData().GetContract()
	for range contractList {
		signature, err := crypto.Sign(hash, privateKey)
		if err != nil {
			s.logger.With(logger.Field{
				"submodule": "sign transaction",
				"error":     err,
			}).ErrorF("error in sign transaction")
			return nil, err
		}
		tx.Signature = append(tx.Signature, signature)
	}

	return json.Marshal(tx)
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.TronWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
	}
	return s
}

func pathFromUserID(uid uint64, index uint32) string {
	str := strconv.FormatUint(uid, 10)
	return fmt.Sprintf("44'/195'/%s'/0/%d", str, index)
}
