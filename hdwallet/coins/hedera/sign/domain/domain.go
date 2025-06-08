package domain

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type Signer interface {
	// SignTransaction signs a Hedera transaction with the given private key and returns the signed bytes or an error.
	SignTransaction(tx *hedera.Transaction, privKey hedera.PrivateKey) ([]byte, error)
}

// TODO: Define Hedera signing domain interfaces here 