package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "pull all entities from server and rewrite local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.Sync(ctx); err != nil {
			fmt.Printf("synchronization of data failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("synchronization of data completed")
	},
}
