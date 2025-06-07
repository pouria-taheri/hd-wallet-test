package tx

import (
	"encoding/hex"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/address"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/crypto"
	"github.com/fxamacker/cbor/v2"
	"golang.org/x/crypto/blake2b"
)

const maxUint64 uint64 = 1<<64 - 1

type txBuilderInput struct {
	Input  transactionInput
	Amount uint64
}

type txBuilderOutput struct {
	Address address.Address
	Amount  uint64
}

type TxBuilder struct {
	Tx       Transaction                               `json:"tx"`
	Protocol ProtocolParams                            `json:"protocol"`
	Inputs   []txBuilderInput                          `json:"inputs"`
	Outputs  []txBuilderOutput                         `json:"outputs"`
	Ttl      uint64                                    `json:"ttl"`
	Fee      uint64                                    `json:"fee"`
	Skeys    map[int]crypto.ExtendedSigningKey         `json:"skeys"`
	Vkeys    map[string]crypto.ExtendedVerificationKey `json:"vkeys"`
	Pkeys    map[string]crypto.ExtendedSigningKey      `json:"pkeys"`
}

func (builder *TxBuilder) Sign(xsk crypto.ExtendedSigningKey) {
	pkeyHashBytes := blake2b.Sum256(xsk)
	pkeyHashString := hex.EncodeToString(pkeyHashBytes[:])
	builder.Pkeys[pkeyHashString] = xsk
}

type ProtocolParams struct {
	MinimumUtxoValue uint64
	PoolDeposit      uint64
	KeyDeposit       uint64
	MinFeeA          uint64
	MinFeeB          uint64
}

type TransactionID string

func (id TransactionID) Bytes() []byte {
	bytes, err := hex.DecodeString(string(id))
	if err != nil {
		panic(err)
	}

	return bytes
}

type Transaction struct {
	_          struct{} `cbor:",toarray"`
	Body       transactionBody
	WitnessSet transactionWitnessSet
	Metadata   *transactionMetadata // or null
}

func (tx *Transaction) Bytes() []byte {
	bytes, err := cbor.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (tx *Transaction) CborHex() string {
	return hex.EncodeToString(tx.Bytes())
}

func (tx *Transaction) ID() TransactionID {
	txHash := blake2b.Sum256(tx.Body.Bytes())
	return TransactionID(hex.EncodeToString(txHash[:]))
}

type transactionWitnessSet struct {
	VKeyWitnessSet []vkeyWitness `cbor:"0,keyasint,omitempty"`
	// TODO: add optional fields 1-4
}

type vkeyWitness struct {
	_         struct{} `cbor:",toarray"`
	VKey      []byte   // ed25519 public key
	Signature []byte   // ed25519 signature
}

// Cbor map
type transactionMetadata map[uint64]transactionMetadatum

// This could be cbor map, array, int, bytes or a text
type transactionMetadatum struct{}

type transactionBody struct {
	Inputs       []transactionInput  `cbor:"0,keyasint"`
	Outputs      []transactionOutput `cbor:"1,keyasint"`
	Fee          uint64              `cbor:"2,keyasint"`
	Ttl          uint64              `cbor:"3,keyasint"`
	Certificates []certificate       `cbor:"4,keyasint,omitempty"` // Omit for now
	Withdrawals  *uint               `cbor:"5,keyasint,omitempty"` // Omit for now
	Update       *uint               `cbor:"6,keyasint,omitempty"` // Omit for now
	MetadataHash *uint               `cbor:"7,keyasint,omitempty"` // Omit for now
}

func (body *transactionBody) Bytes() []byte {
	bytes, err := cbor.Marshal(body)
	if err != nil {
		panic(err)
	}
	return bytes
}

type transactionInput struct {
	_     struct{} `cbor:",toarray"`
	ID    []byte   // HashKey 32 bytes
	Index uint64
}

type transactionOutput struct {
	_       struct{} `cbor:",toarray"`
	Address []byte
	Amount  uint64
}

type certificate struct{}
