# Blockchain
Blockchain and cryptocurrency implementation in Go

## Overview of the project structure:
```bash
├── Cryptography
│   ├── xcrypt.go
│   └── security.go
├── Database
│   ├── block.go
│   ├── database.go
│   ├── genesis.go
│   ├── state.go
│   └── transaction.go
├── Networking
│   ├── client.go
│   ├── http.go
│   └── listener.go
├── Node
│   ├── node.go
│   └── sync.go
└── UI
    └── CommandLine
        ├── cmd
        │   ├── root.go
        │   └── test.go
        ├── main.go
        ├── go.mod
        └── go.sum
```

## How to use HTTP 
Note: Test must be running
For viewing balances 
curl -X GET http://localhost:8080/balances/list 

For adding balances
curl -X POST http://localhost:8080/transaction/create -H "Content-Type: application/json" -d '{"From":"NAME HERE","To":"NAME","Amount":AMOUNT HERE, "Type":"TYPE HERE"}'
