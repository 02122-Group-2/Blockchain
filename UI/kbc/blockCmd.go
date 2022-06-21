package main

import (
	Database "blockchain/Database"
	"fmt"

	"github.com/spf13/cobra"
)

// * Asger, s204435

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

		block := state_block.CreateBlock(state_block.TxMempool)

		err := state_block.AddBlock(block)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Succesfully added the block to the blockchain")
		}
	},
}
