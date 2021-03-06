# KiloBitCoin - Blockchain & Cryptocurrency
Blockchain and cryptocurrency implementation in Go.
Based on the Go implementation of Ethereum, implemented in a simplified version from the bottom up.

Currently, there is no mining algorithm, but everything else, including a consensus algorithm for a Proof-Of-Work type mining system and signing of transactions, is implemented in this proof of concept, developed for educational purposes.

## Overview of the project structure:
```bash
βββ πCryptography // signing of transactions, creation of wallets
β   βββ xcrypt.go
β   βββ wallet_test.go
β   βββ wallet.go
βββ πDatabase // Persistence files and methods for interacting with the blockchain and local node state
β   βββ πPersistence
β   β   βββ πtest_data
β   β   β   βββ ...
β   β   β   βββ test-files
β   β   βββ Blockchain.db
β   β   βββ CurrentState.json
β   β   βββ LatestSnapshot.json
β   β   βββ Transactions.json
β   β   βββ state.json
β   β   βββ PeerSet.json
β   βββ block.go
β   βββ block_test.go
β   βββ state.go
β   βββ state_test.go
β   βββ transaction.go
β   βββ transaction_test.go
βββ πNode // communication with the node network and consensus
β   βββ consensus.go
β   βββ consensus_test.go
β   βββ dtos.go
β   βββ node.go
β   βββ node_test.go
β   βββ receiver.go
β   βββ sender.go
β   βββ set.go
β   βββ sync.go
β   βββ sync_test.go
βββ πShared // utility functions used in other packages
β   βββ *constants.go // must create this yourself
β   βββ fs.go
β   βββ fs_test.go
β   βββ globalconstants.go
β   βββ util.go
βββπUI // command line structure and runtime
   βββ πkbc
       βββ balances.go
       βββ blockCmd.go
       βββ main.go
       βββ overviewCmd.go
       βββ peerCmd.go
       βββ runCmd.go
       βββ transactionCmd.go
       βββ walletCmd.go
```
