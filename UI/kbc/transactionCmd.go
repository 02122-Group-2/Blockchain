package main

import (
	// "fmt"
	// "os"

	// Database "blockchain/Database"

	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagAmount = "amount"
const flagType = "type"

var isCreated = false

func transactionCmd() *cobra.Command {
	var transactionCmd = &cobra.Command{
		Use:   "transaction",
		Short: "Transaction commands (create, load)",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	transactionCmd.AddCommand(transactionCreateCmd())

	return transactionCmd

}

func transactionCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create transaction to database",
		Run: func(cmd *cobra.Command, args []string) {
			//Initialize the flags
			/*
				from, _ := cmd.Flags().GetString(flagFrom)
				to, _ := cmd.Flags().GetString(flagTo)
				amount, _ := cmd.Flags().GetUint(flagAmount)
				typeT, _ := cmd.Flags().GetString(flagType)

				//get the current state
				state := Database.LoadState()
				var transaction Database.SignedTransaction

				//Determine type of transaction
				switch typeT {
				case "genesis":
					transaction = state.CreateGenesisTransaction(Database.AccountAddress(from), float64(amount))

					fmt.Println("Genesis created" + Database.TxToString(transaction))

					isCreated = true

				case "reward":
					transaction = state.CreateReward(Database.AccountAddress(from), float64(amount))

					fmt.Println("Reward created" + Database.TxToString(transaction))

					isCreated = true

				case "transaction":
					if to != "" {
						transaction = state.CreateTransaction(Database.AccountAddress(from), Database.AccountAddress(to), float64(amount))

						fmt.Println("Transaction created" + Database.TxToString(transaction))
						isCreated = true
					}
				}

				if isCreated {
					//Add the transaction to the state and save the transactions
					err := state.AddTransaction(transaction)
					fmt.Println("Transaction succsesfully added")

					Database.SaveTransaction(state.TxMempool)

					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					fmt.Println("TX successfully saved")

				}*/

		},
	}
	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")

	cmd.Flags().Uint(flagAmount, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagAmount)

	cmd.Flags().String(flagType, "", "What type of transaction")
	cmd.MarkFlagRequired(flagType)

	return cmd

}
