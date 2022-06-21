package main

import (
	"fmt"
	"os"

	Crypto "blockchain/Cryptography"
	Database "blockchain/Database"

	"github.com/spf13/cobra"
)

// * Magnus, s204509

func walletCmd() *cobra.Command {
	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Wallet commands (create, delete, access)",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	walletCmd.AddCommand(walletCreateCmd())
	walletCmd.AddCommand(walletDelete())
	walletCmd.AddCommand(walletAccess())

	return walletCmd

}

func walletCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new wallet",
		Run: func(cmd *cobra.Command, args []string) {
			//Initialize the flags

			username, _ := cmd.Flags().GetString(flagUsername)
			password, _ := cmd.Flags().GetString(flagPassword)

			// Create Wallet
			walletAddr, err := Crypto.CreateNewWallet(username, password)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("Wallet has been created with the public address of: " + walletAddr)

		},
	}
	cmd.Flags().String(flagUsername, "", "Username of wallet account")
	cmd.MarkFlagRequired(flagUsername)

	cmd.Flags().String(flagPassword, "", "Password of wallet account")
	cmd.MarkFlagRequired(flagPassword)

	return cmd

}

func walletDelete() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a wallet",
		Run: func(cmd *cobra.Command, args []string) {
			//Initialize the flags

			username, _ := cmd.Flags().GetString(flagUsername)
			password, _ := cmd.Flags().GetString(flagPassword)

			// Access Wallet
			wallet, err := Crypto.AccessWallet(username, password)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Delete Wallet
			err = wallet.Delete(username, password)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("Wallet has been deleted")
		},
	}
	cmd.Flags().String(flagUsername, "", "Username of wallet account")
	cmd.MarkFlagRequired(flagUsername)

	cmd.Flags().String(flagPassword, "", "Password of wallet account")
	cmd.MarkFlagRequired(flagPassword)

	return cmd

}

func walletAccess() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "access",
		Short: "Aceess the wallet to see its information",
		Run: func(cmd *cobra.Command, args []string) {
			//Initialize the flags

			username, _ := cmd.Flags().GetString(flagUsername)
			password, _ := cmd.Flags().GetString(flagPassword)

			// Access Wallet
			wallet, err := Crypto.AccessWallet(username, password)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			//get the current state
			state := Database.LoadState()

			// Get private key
			privKey, _ := wallet.GetPrivateKey(password)

			// List its information
			fmt.Println("Address of wallet: " + wallet.Address)
			fmt.Println("Private key of wallet: " + privKey.X.String())
			fmt.Println("Balance of wallet:" + fmt.Sprint(state.AccountBalances[Database.AccountAddress(wallet.Address)]))

		},
	}
	cmd.Flags().String(flagUsername, "", "Username of wallet account")
	cmd.MarkFlagRequired(flagUsername)

	cmd.Flags().String(flagPassword, "", "Password of wallet account")
	cmd.MarkFlagRequired(flagPassword)

	return cmd

}
