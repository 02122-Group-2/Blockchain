package main

import (
	Database "blockchain/Database"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

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

			fmt.Printf("Latest hash: %v \n", currentState.LatestHash)

			//converting time to readable format
			tUnix := currentState.LastBlockTimestamp / int64(time.Second)
			tUnixNanoRemainder := (currentState.LastBlockTimestamp % int64(time.Second))
			formattedTime := time.Unix(tUnix, tUnixNanoRemainder)

			fmt.Printf("Latest Block Timestamp: %v", formattedTime)
		},
	}

	return overviewCmd
}
