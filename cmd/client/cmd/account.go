/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
	"golang.org/x/term"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "create or login account",
	Long:  `create or login account on server for getting access for work with private data`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create account on server",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fmt.Print("password: ")
		bytePwd, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("failed to read password")
			os.Exit(reasonInternalError)
		}
		fmt.Print("\n")

		os.Exit(register(ctx, app, string(bytePwd), args))
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login account on server",
	Long:  `login account and save credential on workspace folder`,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fmt.Print("password: ")
		bytePwd, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("failed to read password")
			os.Exit(reasonInternalError)
		}
		fmt.Print("\n")

		os.Exit(login(ctx, app, string(bytePwd), args))
	},
}

func init() {
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(loginCmd)
}

func login(ctx context.Context, app client.Client, readPassword string, args []string) int {
	if len(args) == 0 {
		fmt.Println("does not define account name")
		return reasonNotFillRequiredArgs
	}
	accountName := args[0]

	if err := app.Login(ctx, accountName, readPassword); err != nil {
		fmt.Println("login failed")
		fmt.Println(err.Error())
		return reasonInternalError
	}
	fmt.Println("login operation completed")
	return 0
}

func register(ctx context.Context, app client.Client, readPassword string, args []string) int {
	if len(args) == 0 {
		fmt.Println("does not define account name")
		return reasonNotFillRequiredArgs
	}
	accountName := args[0]

	if err := app.CreateAccount(ctx, accountName, readPassword); err != nil {
		fmt.Println("failed to create account")
		fmt.Println(err.Error())
		return reasonInternalError
	}
	fmt.Println("account created")
	return 0
}
