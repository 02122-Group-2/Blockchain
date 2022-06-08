package main

import (
	"fmt"
	"os"

	Node "blockchain/Node"

	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches the node and HTTP API",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Starting up...")

			err := Node.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	//addDefaultRequiredFlags(runCmd)

	return runCmd

}

