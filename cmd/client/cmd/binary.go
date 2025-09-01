package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
)

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "manager of binary data",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type binaryAddFlags struct {
	file string
	name string
}

var binaryAdd = binaryAddFlags{}

var binaryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a entity from file and save it on server and on local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		os.Exit(binaryAddHandle(ctx, app))
	},
}

type binaryGetFlags struct {
	id int
}

var binaryGet = binaryGetFlags{}

var binaryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a entity from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		os.Exit(binaryGetHandle(ctx, app))
	},
}

type binaryEditFlags struct {
	id   int
	name string
	file string
}

var binaryEdit = binaryEditFlags{}

var binaryEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a entity on server and on local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		os.Exit(binaryEditHandle(ctx, app))

	},
}

type binaryDeleteFlags struct {
	id int
}

var binaryDelete = binaryDeleteFlags{}

var binaryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a entity on server and on local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()
		// TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		os.Exit(binaryDeleteHandle(ctx, app))
	},
}

func init() {
	binaryAddCmd.PersistentFlags().StringVarP(&binaryAdd.file, "file", "", "", "file")
	binaryAddCmd.PersistentFlags().StringVarP(&binaryAdd.name, "name", "", "", "name")

	binaryGetCmd.PersistentFlags().IntVarP(&binaryGet.id, "id", "", 0, "id")

	binaryEditCmd.PersistentFlags().IntVarP(&binaryEdit.id, "id", "", 0, "id")
	binaryEditCmd.PersistentFlags().StringVarP(&binaryEdit.file, "file", "", "", "file")
	binaryEditCmd.PersistentFlags().StringVarP(&binaryEdit.name, "name", "", "", "name")

	binaryDeleteCmd.PersistentFlags().IntVarP(&binaryDelete.id, "id", "", 0, "id")

	binaryCmd.AddCommand(binaryAddCmd)
	binaryCmd.AddCommand(binaryGetCmd)
	binaryCmd.AddCommand(binaryEditCmd)
	binaryCmd.AddCommand(binaryDeleteCmd)
}

func binaryAddHandle(ctx context.Context, app client.Client) int {
	if binaryAdd.file == "" {
		fmt.Println("--file argument is required")
		return reasonNotFillRequiredArgs
	}

	if binaryAdd.name == "" {
		fmt.Println("--name argument is required")
		return reasonNotFillRequiredArgs
	}

	metadata := prepareMetadata()

	content, err := os.ReadFile(binaryAdd.file)
	if err != nil {
		fmt.Printf("reading file failed: %s\n", err)
		return reasonInternalError
	}

	if err := app.SaveBinaryData(ctx, content, binaryAdd.name, 0, metadata); err != nil {
		fmt.Printf("saving binary data failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("saving binary data success")
	return 0
}

func binaryEditHandle(ctx context.Context, app client.Client) int {
	if binaryEdit.file == "" {
		fmt.Println("--file argument is required")
		return reasonNotFillRequiredArgs
	}

	if binaryEdit.name == "" {
		fmt.Println("--name argument is required")
		return reasonNotFillRequiredArgs
	}

	if binaryEdit.id == 0 {
		fmt.Println("--id argument is required")
		return reasonNotFillRequiredArgs
	}
	metadata := prepareMetadata()

	content, err := os.ReadFile(binaryEdit.file)
	if err != nil {
		fmt.Printf("reading file failed: %s\n", err)
		return reasonInternalError
	}

	if err := app.SaveBinaryData(ctx, content, binaryEdit.name, binaryEdit.id, metadata); err != nil {
		fmt.Printf("saving text failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("saving binary data success")
	return 0
}

func binaryGetHandle(ctx context.Context, app client.Client) int {
	if err := app.GetBinaryData(ctx, binaryGet.id); err != nil {
		fmt.Printf("getting binary data failed: %s\n", err)
		return reasonInternalError
	}
	return 0
}

func binaryDeleteHandle(ctx context.Context, app client.Client) int {
	if binaryDelete.id == 0 {
		fmt.Println("--id argument is required")
		return reasonNotFillRequiredArgs
	}

	if err := app.DeleteBinaryData(ctx, binaryDelete.id); err != nil {
		fmt.Printf("deleting binary data failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("deleting binary data is completed")
	return 0
}
