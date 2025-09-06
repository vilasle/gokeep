package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "manager for work with text data",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type textAddFlags struct {
	data string
	name string
}

var textAdd = textAddFlags{}

var textAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add new entity",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := textAddHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

type textGetFlags struct {
	id int
}

var textGet = textGetFlags{}

var textGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get all entities or specific entity(use --id argument)",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := textGetHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

type textEditFlags struct {
	id   int
	data string
	name string
}

var textEdit = textEditFlags{}

var textEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit existed entity",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := textEditHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

type textDeleteFlags struct {
	id int
}

var textDelete = textDeleteFlags{}

var textDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete entity from server and local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := textDeleteHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

func init() {
	textAddCmd.PersistentFlags().StringVarP(&textAdd.data, "data", "", "", "text sting which need to save")
	textAddCmd.PersistentFlags().StringVarP(&textAdd.name, "name", "", "", "name of entity")

	textGetCmd.PersistentFlags().IntVarP(&textGet.id, "id", "", 0, "id")

	textEditCmd.PersistentFlags().IntVarP(&textEdit.id, "id", "", 0, "id")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.data, "data", "", "", "text sting which need to save")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.name, "name", "", "", "name of entity")

	textDeleteCmd.PersistentFlags().IntVarP(&textDelete.id, "id", "", 0, "id")

	textCmd.AddCommand(textAddCmd)
	textCmd.AddCommand(textGetCmd)
	textCmd.AddCommand(textEditCmd)
	textCmd.AddCommand(textDeleteCmd)
}

func textAddHandle(ctx context.Context, app client.Client) int {
	if textAdd.name == "" {
		fmt.Println("--name argument is required")
		return reasonNotFillRequiredArgs
	}

	if textAdd.data == "" {
		fmt.Println("'--data' argument is required")
		return reasonNotFillRequiredArgs
	}

	metadata := prepareMetadata()

	if err := app.SaveTextData(ctx, textAdd.data, textAdd.name, 0, metadata); err != nil {
		fmt.Printf("saving text failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("saving text success")

	return 0
}

func textEditHandle(ctx context.Context, app client.Client) int {
	if textEdit.name == "" {
		fmt.Println("--name argument is required")
		return reasonNotFillRequiredArgs
	}

	if textEdit.data == "" {
		fmt.Println("'data' argument is required")
		return reasonNotFillRequiredArgs
	}

	if textEdit.id == 0 {
		fmt.Println("--id is required")
		return reasonNotFillRequiredArgs
	}

	metadata := prepareMetadata()

	if textEdit.data != "" {
		if err := app.SaveTextData(ctx, textEdit.data, textEdit.name, textEdit.id, metadata); err != nil {
			fmt.Printf("saving text failed: %s\n", err)
			return reasonInternalError
		}
		fmt.Println("saving text success")
	}
	return 0
}

func textGetHandle(ctx context.Context, app client.Client) int {
	if err := app.GetTextData(ctx, textGet.id); err != nil {
		fmt.Printf("getting text data failed: %s\n", err)
		return reasonInternalError
	}
	return 0
}

func textDeleteHandle(ctx context.Context, app client.Client) int {
	if textDelete.id == 0 {
		fmt.Println("--id argument is required")
		return reasonNotFillRequiredArgs
	}

	if err := app.DeleteTextData(ctx, textDelete.id); err != nil {
		fmt.Printf("deleting text data failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("deleting text data is completed")
	return 0
}
