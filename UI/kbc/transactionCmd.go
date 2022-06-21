package main

import (
	"fmt"
	"os"

	Crypto "blockchain/Cryptography"
	Database "blockchain/Database"

	"github.com/spf13/cobra"
)

// * Asger, s204435

const flagUsername = "username"
const flagPassword = "password"
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
			//Initialize the flags

			username, _ := cmd.Flags().GetString(flagUsername)
			password, _ := cmd.Flags().GetString(flagPassword)
			toRaw, _ := cmd.Flags().GetString(flagTo)
			to := Database.AccountAddress(toRaw)
			amount, _ := cmd.Flags().GetUint(flagAmount)

			// Access Wallet
			wallet, err := Crypto.AccessWallet(username, password)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			//get the current state
			state := Database.LoadState()
			var transaction Database.SignedTransaction

			//Determine type of transaction
			if to != "" {

				transaction, err = state.CreateSignedTransaction(wallet, password, to, float64(amount))

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Println("Transaction created: " + Database.TxToString(transaction.Tx))

				err = state.AddTransaction(transaction)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				fmt.Println("Transaction succsesfully saved")

			} else {
				fmt.Fprintln(os.Stderr, "Sender is undefined")
				os.Exit(1)
			}

		},
	}
	cmd.Flags().String(flagUsername, "", "Username of wallet account")
	cmd.MarkFlagRequired(flagUsername)

	cmd.Flags().String(flagPassword, "", "Password of wallet account")
	cmd.MarkFlagRequired(flagPassword)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagAmount, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagAmount)

	return cmd

}
