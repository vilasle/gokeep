/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "",
	Long:  ``,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called account create")
		if len(args) == 0 {
			fmt.Println("does not define account name")
			os.Exit(1)
		}
		accountName := args[0]
		fmt.Println("account name:", accountName)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called account login")
		if len(args) == 0 {
			fmt.Println("does not define account name")
			os.Exit(1)
		}
		accountName := args[0]
		fmt.Println("account name:", accountName)
	},
}

func init() {
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(loginCmd)
}
