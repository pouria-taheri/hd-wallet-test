module git.mazdax.tech/blockchain/hdwallet

go 1.14

require (
	filippo.io/edwards25519 v1.0.0-rc.1
	git.mazdax.tech/blockchain/bnb-go-sdk v1.2.15
	git.mazdax.tech/core/errors v0.1.3
	git.mazdax.tech/core/swag v1.8.1
	git.mazdax.tech/core/swaggerfiles v0.0.0-20210629070638-edad3818fe8f
	git.mazdax.tech/data-layer/configcore v0.0.0-20210629062014-3c80d2c26989
	git.mazdax.tech/data-layer/loggercore v0.1.0
	git.mazdax.tech/data-layer/viper v0.1.1
	git.mazdax.tech/delivery/ginger v0.1.5
	git.mazdax.tech/delivery/handlercore v0.1.4
	github.com/aliworkshop/stellar-go v0.1.2
	github.com/amintalebi/go-subkey v1.0.4
	github.com/bchsuite/bchd v0.0.0-20171115031648-585afbdef8b8
	github.com/bchsuite/bchlog v0.0.0-20171105124706-ed5d9caa3ea3
	github.com/bchsuite/bchutil v0.0.0-20171105124755-1db7414eadd2
	github.com/binance-chain/go-sdk v1.2.6 // indirect
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/btcutil v1.0.2
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/echovl/bech32 v0.1.0
	github.com/echovl/ed25519 v0.2.0
	github.com/eoscanada/eos-go v0.10.0
	github.com/ethereum/go-ethereum v1.10.14
	github.com/etherlabsio/healthcheck v0.0.0-20191224061800-dd3d2fd8c3f6
	github.com/fatih/structs v1.1.0
	github.com/fbsobreira/gotron-sdk v0.0.0-20201228180255-4c6c1768cd2a
	github.com/fxamacker/cbor/v2 v2.2.0
	github.com/gagliardetto/solana-go v1.0.2
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.5.2
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/islishude/bip32 v1.0.2
	github.com/ltcsuite/ltcd v0.20.1-beta
	github.com/ltcsuite/ltcutil v0.0.0-20191227053721-6bec450ea6ad
	github.com/matoous/go-nanoid/v2 v2.0.0
	github.com/pborman/uuid v1.2.1
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.1
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/xana-rahmani/ripple v0.0.1
	github.com/zondax/ledger-go v0.12.2 // indirect
	github.com/hashgraph/hedera-sdk-go v2.28.0 // or latest stable version
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/zondax/ledger-go v0.12.2 => github.com/binance-chain/ledger-go v0.9.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
