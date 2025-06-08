package usecase

import (
	"testing"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func TestSignTransaction(t *testing.T) {
	// Create a dummy logger (replace with a real or mock logger as needed)
	var dummyLogger logger.Logger
	uc := &useCase{logger: dummyLogger}

	// Generate a new private key for testing
	privKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Create a simple transaction (e.g., AccountCreateTransaction)
	tx := hedera.NewAccountCreateTransaction().SetKey(privKey.PublicKey())

	// Sign the transaction
	signedBytes, err := uc.SignTransaction(&tx.Transaction, privKey)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}
	if len(signedBytes) == 0 {
		t.Error("Signed transaction bytes should not be empty")
	}
} 