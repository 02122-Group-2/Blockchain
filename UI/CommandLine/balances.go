package main

import (
	"fmt"
	"os"

	Database "github.com/blockchainProject/blockchain/Database"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",

		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	balancesCmd.AddCommand(balancesListCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		//Load balances here
		state, err := Database.LoadState()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		//Printing balances
		fmt.Println("Accounts balances:")
		fmt.Println("__________________")
		fmt.Println("")
		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s:%d", account, balance))
		}

	},
}
