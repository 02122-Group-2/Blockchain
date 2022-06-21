# KiloBitCoin - Blockchain & Cryptocurrency
Blockchain and cryptocurrency implementation in Go.
Based on the Go implementation of Ethereum, implemented in a simplified version from the bottom up.

Currently, there is no mining algorithm, but everything else, including a consensus algorithm for a Proof-Of-Work type mining system and signing of transactions, is implemented in this proof of concept, developed for educational purposes.

## Overview of the project structure:
```bash
├── 📁Cryptography // signing of transactions, creation of wallets
│   ├── xcrypt.go
│   ├── wallet_test.go
│   └── wallet.go
├── 📁Database // Persistence files and methods for interacting with the blockchain and local node state
│   ├── 📁Persistence
│   │   ├── 📁test_data
│   │   │   ├── ...
│   │   │   └── test-files
│   │   ├── Blockchain.db
│   │   ├── CurrentState.json
│   │   ├── LatestSnapshot.json
│   │   ├── Transactions.json
│   │   ├── state.json
│   │   └── PeerSet.json
│   ├── block.go
│   ├── block_test.go
│   ├── state.go
│   ├── state_test.go
│   ├── transaction.go
│   └── transaction_test.go
├── 📁Node // communication with the node network and consensus
│   ├── consensus.go
│   ├── consensus_test.go
│   ├── dtos.go
│   ├── node.go
│   ├── node_test.go
│   ├── receiver.go
│   ├── sender.go
│   ├── set.go
│   ├── sync.go
│   └── sync_test.go
├── 📁Shared // utility functions used in other packages
│   ├── *constants.go // must create this yourself
│   ├── fs.go
│   ├── fs_test.go
│   ├── globalconstants.go
│   └── util.go
└──📁UI // command line structure and runtime
   └── 📁kbc
       ├── balances.go
       ├── blockCmd.go
       ├── main.go
       ├── overviewCmd.go
       ├── peerCmd.go
       ├── runCmd.go
       ├── transactionCmd.go
       └── walletCmd.go
```
