module blockchain

replace database => ./Database

go 1.17

require github.com/spf13/cobra v1.3.0

require (
	github.com/btcsuite/btcd/btcec/v2 v2.1.2 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/rjeczalik/notify v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)

require (
	github.com/ethereum/go-ethereum v1.10.17
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
