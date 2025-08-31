/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
	"github.com/vilasle/gokeep/internal/client/cli"
)

const (
	reasonWrongConfig         = 1
	reasonNotFillRequiredArgs = 2
	reasonInternalError       = 3
)

var (
	customWorkspace string
	Version         string
	Date            string
	Commit          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gokeep",
	Short: "",
	Long:  ``,
}

// versionCmd represents the 'config report' command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version - print build inforation",
	Long:  `version - print information about build date, version, and commit hash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", Version)
		fmt.Println("Date:", Date)
		fmt.Println("Commit:", Commit)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&customWorkspace,
		"workspace", "c", "",
		"custom workplace, with gokeep directories and files. Default workplace is $HOME/.gokeep")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(dataCmd)
}

func initCLIClient() client.Client {
	config, err := client.GetCurrentConfiguration(customWorkspace)
	if err != nil {
		fmt.Println("getting current configuration failed:")
		fmt.Println(err)
		fmt.Println("\nif you did not initialize configuration try call before creating account:")
		fmt.Println("\tgokeep config init --grpc-socket $GRPC_SOCKET --db-path $DB_PATH")
		os.Exit(reasonWrongConfig)
	}

	app, err := cli.NewClient(config)
	if err != nil {
		fmt.Println("failed to create client")
		fmt.Println(err)
		os.Exit(reasonInternalError)
	}

	return app
}
