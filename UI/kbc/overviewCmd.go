package main

import (
	Database "blockchain/Database"
	Node "blockchain/Node"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// * Emilie, s204471

func overviewCmd() *cobra.Command {
	var overviewCmd = &cobra.Command{
		Use:   "overview",
		Short: "View an overview of the blockchain",

		Run: func(cmd *cobra.Command, args []string) {
			//Print user alias here

			//Print account balance

			//most recent block info
			currentState := Database.LoadState()
			numberOfBlocks := currentState.LastBlockSerialNo
			fmt.Printf("Current number of blocks in blockchain: %v \n", numberOfBlocks)

			fmt.Printf("Latest hash: %x \n", currentState.LatestHash)

			//converting time to readable format
			tUnix := currentState.LastBlockTimestamp / int64(time.Second)
			tUnixNanoRemainder := (currentState.LastBlockTimestamp % int64(time.Second))
			formattedTime := time.Unix(tUnix, tUnixNanoRemainder)

			fmt.Printf("Latest Block Timestamp: %v \n", formattedTime)

			//Number of actve peers
			currentPeers := Node.GetNode().PeerSet

			//Counting all active peers
			var numActivePeers = 0
			for peer, _ := range currentPeers {
				if !Node.Ping(peer).Ok {
					continue
				} else {
					numActivePeers += 1
				}
			}

			fmt.Printf("Number of active peers: %v", numActivePeers)

		},
	}

	return overviewCmd
}
