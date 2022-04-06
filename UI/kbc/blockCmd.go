package main

import (
	Database "blockchain/database"
	"fmt"

	"github.com/spf13/cobra"
)

func blockCmd() *cobra.Command {
	var blockCmd = &cobra.Command{
		Use:   "block",
		Short: "Interact with blocks",

		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	blockCmd.AddCommand(blockCreateCmd)

	return blockCmd
}

var blockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a block and saves it",
	Run: func(cmd *cobra.Command, args []string) {
		var state_block = Database.LoadState()

		pendingTxs := Database.LoadTransactions()

		block := state_block.CreateBlock(pendingTxs)

		fmt.Println(state_block.AddBlock(block).Error())
	},
}
