package main

import (
	"fmt"
	"os"

	Database "github.com/blockchainProject/blockchain/Database"

	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagAmount = "amount"

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
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			amount, _ := cmd.Flags().GetUint(flagAmount)
			state, _ := Database.LoadState()
			fmt.Printf("%v", amount)

			transaction := state.CreateTransaction(Database.AccountAddress(from), Database.AccountAddress(to), float64(amount))
			fmt.Println("Transaction created" + Database.TxToString(transaction))

			//Get state, add transaction and save new state

			err := state.AddTransaction(transaction)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("TX successfully added")

		},
	}
	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagAmount, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagAmount)

	return cmd

}
