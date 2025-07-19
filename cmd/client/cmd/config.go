package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	grpcSocket        string
	localDatabasePath string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config ",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

// initCmd represents the 'config init' command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create necessary file and directories",
	Long: `init - create main directory with config files
	generate RSA keys and create config file with information about local database and grpc server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called 'config init'")
		fmt.Println("grpc socket:", grpcSocket)
		fmt.Println("local database path:", localDatabasePath)
		fmt.Println("config file", customConfig)
	},
}

// reportCmd represents the 'config report' command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "report - print current config values",
	Long:  `report - check and print current config values`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called 'config report'")
		fmt.Println("config file", customConfig)
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&grpcSocket, "grpc-socket", "s", "", "grpc socket")
	initCmd.PersistentFlags().StringVarP(&localDatabasePath, "db-path", "d", "", "local database path")

	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(reportCmd)

}
