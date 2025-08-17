/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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
		if len(args) == 0 {
			fmt.Println("does not define account name")
			os.Exit(reasonNotFillRequiredArgs)
		}
		accountName := args[0]

		app := initCLIClient()
		defer app.Close()

		fmt.Println("login:", accountName)
		fmt.Print("password: ")
		// bytePwd, err := term.ReadPassword(int(syscall.Stdin))
		// if err != nil {
		// 	fmt.Println("failed to read password")
		// 	os.Exit(reasonInternalError)
		// }
		// fmt.Print("\n")
		// pass := string(bytePwd)

		pass := "some password"

		//TODO cancel if got signal
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := app.CreateAccount(ctx, accountName, pass); err != nil {
			fmt.Println("failed to create account")
			os.Exit(reasonInternalError)
		}
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
			os.Exit(reasonNotFillRequiredArgs)
		}
		accountName := args[0]

		fmt.Println("login:", accountName)
		fmt.Print("password: ")
		// bytePwd, err := term.ReadPassword(int(syscall.Stdin))
		// if err != nil {
		// 	fmt.Println("failed to read password")
		// 	os.Exit(1)
		// }
		// fmt.Print("\n")
		// pass := string(bytePwd)

		pass := "some password"

		app := initCLIClient()
		defer app.Close()

		if err := app.Login(accountName, pass); err != nil {
			fmt.Println("failed to create account")
			os.Exit(reasonInternalError)
		}
	},
}

func init() {
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(loginCmd)
}
