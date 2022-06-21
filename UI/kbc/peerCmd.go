package main

import (
	"fmt"

	node "blockchain/Node"
	shared "blockchain/Shared"

	"github.com/spf13/cobra"
)

// * Magnus, s204509

func peerCmd() *cobra.Command {
	var peerCmd = &cobra.Command{
		Use:   "peer",
		Short: "Peer commands (add)",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	peerCmd.AddCommand(peerAddCommand())

	return peerCmd

}

func peerAddCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add [Address]",
		Short: "Add peer to set of peers",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires the argument of a single address")
			}
			if shared.LegalIpAddress(args[0]) {
				return nil
			}
			return fmt.Errorf("invalid Ip Address: %s", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			//Check valid argument
			addr := args[0]
			peerset := node.GetPeerSet()
			if !node.Ping(addr).Ok {
				fmt.Println("Error: Unable to establish a connection with " + addr)
			} else {
				peerset.Add(addr)
				node.SavePeerSetAsJSON(peerset, shared.PeerSetFile)
				fmt.Println("Succesfully saved the peer address " + addr)
			}

		},
	}
	return cmd

}
