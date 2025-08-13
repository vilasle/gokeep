/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
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
			os.Exit(1)
		}
		accountName := args[0]

		config, err := client.GetCurrentConfiguration(customWorkspace)
		if err != nil {
			fmt.Println("getting current configuration failed:")
			fmt.Println(err)
			fmt.Println("\nif you did not initialize configuration try call before creating account:")
			fmt.Println("\tgokeep config init --grpc-socket $GRPC_SOCKET --db-path $DB_PATH")
			os.Exit(2)
		}

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

		app, err := client.NewClient(config)
		if err != nil {
			fmt.Println("failed to create client")
			os.Exit(1)
		}

		if err := app.CreateAccount(accountName, pass); err != nil {
			fmt.Println("failed to create account")
			os.Exit(1)
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
			os.Exit(1)
		}
		accountName := args[0]
		config, err := client.GetCurrentConfiguration(customWorkspace)
		if err != nil {
			fmt.Println("getting current configuration failed:")
			fmt.Println(err)
			fmt.Println("\nif you did not initialize configuration try call before creating account:")
			fmt.Println("\tgokeep config init --grpc-socket $GRPC_SOCKET --db-path $DB_PATH")
			os.Exit(2)
		}

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

		app, err := client.NewClient(config)
		if err != nil {
			fmt.Println("failed to create client")
			os.Exit(1)
		}

		if err := app.Login(accountName, pass); err != nil {
			fmt.Println("failed to create account")
			os.Exit(1)
		}
	},
}

func init() {
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(loginCmd)
}
