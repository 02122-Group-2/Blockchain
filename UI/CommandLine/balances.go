package main

import (
	"fmt"

	Database "blockchain/database"

	"github.com/spf13/cobra"

	"os"
	"path/filepath"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",

		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	addDefaultRequiredFlags(balancesListCmd)

	balancesCmd.AddCommand(balancesListCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		//Load balances here
		state := Database.LoadState()

		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)

		str, _ := cmd.Flags().GetString(currentDir)
		fmt.Println(str)

		//Printing balances
		fmt.Println("Accounts balances:")
		fmt.Println("__________________")
		fmt.Println("")
		for account, balance := range state.AccountBalances {
			fmt.Printf("%s:%d\n", account, balance)
		}

	},
}
