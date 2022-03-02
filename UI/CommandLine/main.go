package main

import (
	"fmt"
	"os"

	"../../database"
	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "1"
const Fix = "0"
const Verbal = "TX Add && Balances List"

func main() {
	state := database.State.LoadState()
	state.getLatestHash()

	var tbbCmd = &cobra.Command{
		Use:   "Monkeycoin",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Describes version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s.%s.%s-beta %s.", Major, Minor, Fix, Verbal)
		},
	}

	tbbCmd.AddCommand(versionCmd)

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
