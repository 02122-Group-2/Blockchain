package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",

		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	balancesCmd.AddCommand(balancesCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		//Load balances here

		fmt.Println("Accounts balances:")
		fmt.Println("__________________")
		fmt.Println("")
		//Print balances here

	},
}
