package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	var kbcCmd = &cobra.Command{
		Use:   "kbc",
		Short: "The Kaelder Bar Coin CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	kbcCmd.AddCommand(versionCmd)
	kbcCmd.AddCommand(transactionCmd())
	kbcCmd.AddCommand(balancesCmd())
	kbcCmd.AddCommand(blockCmd())

	err := kbcCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const Major = "0"
const Minor = "1"
const Fix = "0"
const Verbal = "TX Add && Balances List"
const flagDataDir = "datadir"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s-beta %s.", Major, Minor, Fix, Verbal)
	},
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")
	cmd.MarkFlagRequired(flagDataDir)

}
