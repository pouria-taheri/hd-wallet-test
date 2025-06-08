package usecase

import (
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"git.mazdax.tech/blockchain/hdwallet/coins/hedera/sign/domain"
)

type useCase struct {
	logger logger.Logger
}

func New(logger logger.Logger, configRegistry configcore.Registry) *useCase {
	return &useCase{
		logger: logger,
	}
}

// Ensure useCase implements domain.Signer
var _ domain.Signer = (*useCase)(nil)

// SignTransaction signs a Hedera transaction with the given private key and returns the signed bytes or an error.
func (u *useCase) SignTransaction(tx *hedera.Transaction, privKey hedera.PrivateKey) ([]byte, error) {
	signedTx, err := tx.Sign(privKey)
	if err != nil {
		u.logger.Error("Failed to sign transaction", err)
		return nil, err
	}
	return signedTx.ToBytes()
}

// TODO: Implement Hedera transaction signing logic here 