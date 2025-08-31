package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "",
	Long:  ``,
}

type binaryAddFlags struct {
	file string
	name string
}

var binaryAdd = binaryAddFlags{}

var binaryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if binaryAdd.file == "" {
			fmt.Println("--file argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if binaryAdd.name == "" {
			fmt.Println("--name argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		metadata := prepareMetadata()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		content, err := os.ReadFile(binaryAdd.file)
		if err != nil {
			fmt.Printf("reading file failed: %s\n", err)
			os.Exit(reasonInternalError)
		}

		if err := app.SaveBinaryData(ctx, content, binaryAdd.name, 0, metadata); err != nil {
			fmt.Printf("saving binary data failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("saving binary data success")
	},
}

type binaryGetFlags struct {
	id int
}

var binaryGet = binaryGetFlags{}

var binaryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.GetBinaryData(ctx, binaryGet.id); err != nil {
			fmt.Printf("getting binary data failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
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
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if binaryEdit.file == "" {
			fmt.Println("--file argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if binaryEdit.name == "" {
			fmt.Println("--name argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if binaryEdit.id == 0 {
			fmt.Println("--id argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}
		metadata := prepareMetadata()
		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		content, err := os.ReadFile(binaryEdit.file)
		if err != nil {
			fmt.Printf("reading file failed: %s\n", err)
			os.Exit(reasonInternalError)
		}

		if err := app.SaveBinaryData(ctx, content, binaryEdit.name, binaryEdit.id, metadata); err != nil {
			fmt.Printf("saving text failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("saving binary data success")
	},
}

type binaryDeleteFlags struct {
	id int
}

var binaryDelete = binaryDeleteFlags{}

var binaryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if textDelete.id == 0 {
			fmt.Println("--id argument is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.DeleteBinaryData(ctx, binaryDelete.id); err != nil {
			fmt.Printf("deleting binary data failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("deleting binary data is completed")
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
