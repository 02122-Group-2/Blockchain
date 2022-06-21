package main

import (
	"fmt"
	"os"

	shared "blockchain/Shared"

	"github.com/spf13/cobra"
)

// * Asger, s204435

func main() {
	var kbcCmd = &cobra.Command{
		Use:   "kbc",
		Short: "The KiloBitCoin CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	shared.EnsureNeededFilesExist()

	kbcCmd.AddCommand(versionCmd)
	kbcCmd.AddCommand(transactionCmd())
	kbcCmd.AddCommand(balancesCmd())
	kbcCmd.AddCommand(runCmd())
	kbcCmd.AddCommand(blockCmd())
	kbcCmd.AddCommand(overviewCmd())
	kbcCmd.AddCommand(walletCmd())
	kbcCmd.AddCommand(peerCmd())

	err := kbcCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const Major = "2"
const Minor = "0"
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
