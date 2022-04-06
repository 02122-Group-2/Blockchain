package main

import (
	"fmt"
	"os"

	Node "blockchain/node"

	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches the node and HTTP API",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			fmt.Println("Starting up...")

			err := Node.Run(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	addDefaultRequiredFlags(runCmd)

	return runCmd

}
