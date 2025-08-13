package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
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
		if grpcSocket == "" {
			fmt.Println("grpc socket is not set")
			os.Exit(2)
		}

		if localDatabasePath == "" {
			fmt.Println("local database path is not set")
			os.Exit(3)
		}

		err := client.CreateNewConfiguration(customWorkspace, grpcSocket, localDatabasePath)
		if err != nil {
			fmt.Println("creating new configuration failed:", err)
			os.Exit(4)
		}
		fmt.Println("configuration created")
	},
}

// reportCmd represents the 'config report' command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "report - print current config values",
	Long:  `report - check and print current config values`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.GetCurrentConfiguration(customWorkspace)
		if err != nil {
			fmt.Println("getting current configuration failed:", err)
			os.Exit(2)
		}
		fmt.Println("current configuration:")
		config.Report()
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&grpcSocket, "grpc-socket", "s", "", "grpc socket")
	initCmd.PersistentFlags().StringVarP(&localDatabasePath, "db-path", "d", "", "local database path")

	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(reportCmd)

}
