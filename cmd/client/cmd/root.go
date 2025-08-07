/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	customWorkplace string
	Version      string
	Date         string
	Commit       string
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
	rootCmd.PersistentFlags().StringVarP(&customWorkplace,
		"workplace", "c", "",
		"custom workplace, with gokeep directories and files. Default workplace is $HOME/.gokeep")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(dataCmd)
}
