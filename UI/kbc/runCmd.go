package main

import (
	"fmt"
	"os"

	Node "blockchain/Node"
	shared "blockchain/Shared"

	"github.com/spf13/cobra"
)

// * Emilie, s204471

var setPort = "port"
var setBootstrap = "bootstrap"

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches the node and HTTP API",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Starting up...")

			port, _ := cmd.Flags().GetInt(setPort)
			bootstrp, _ := cmd.Flags().GetString(setBootstrap)

			fmt.Println(shared.BootstrapNode)
			fmt.Println(bootstrp)

			shared.HttpPort = port
			shared.BootstrapNode = bootstrp

			fmt.Println(shared.BootstrapNode)
			fmt.Println(bootstrp)

			err := Node.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	runCmd.Flags().Int(setPort, shared.HttpPort, "Optinal Flag: Manually set port. Default is "+string(shared.HttpPort))
	runCmd.Flags().String(setBootstrap, shared.BootstrapNode, "Optional Flag: Manually set bootstrap node. Default is "+shared.BootstrapNode)
	//addDefaultRequiredFlags(runCmd)

	return runCmd

}
