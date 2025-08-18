package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "",
	Long:  ``,
}

type textAddFlags struct {
	data string
	file string
	name string
}

var textAdd = textAddFlags{}

var textAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if textAdd.name == "" {
			fmt.Println("name is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if textAdd.data == "" && textAdd.file == "" {
			fmt.Println("'data' or 'file' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if textAdd.data != "" {
			if err := app.SaveTextDataAsIs(ctx, textAdd.data, textAdd.name, 0); err != nil {
				fmt.Printf("saving text failed: %s\n", err)
				os.Exit(reasonInternalError)
			}
			fmt.Println("saving text success")
		} else {
			if err := app.SaveTextDataFromFile(ctx, textAdd.file, textAdd.name, 0); err != nil {
				fmt.Printf("saving text failed: %s\n", err)
				os.Exit(reasonInternalError)
			}
			fmt.Println("saving text success")
		}
	},
}

type textGetFlags struct {
	id int
}

var textGet = textGetFlags{}

var textGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.GetTextData(ctx, textGet.id); err != nil {
			fmt.Printf("getting text data failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
	},
}

type textEditFlags struct {
	id   int
	data string
	file string
	name string
}

var textEdit = textEditFlags{}

var textEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if textEdit.name == "" {
			fmt.Println("name is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if textEdit.data == "" && textEdit.file == "" {
			fmt.Println("'data' or 'file' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if textEdit.id == 0 {
			fmt.Println("id is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if textAdd.data != "" {
			if err := app.SaveTextDataAsIs(ctx, textEdit.data, textEdit.name, textEdit.id); err != nil {
				fmt.Printf("saving text failed: %s\n", err)
				os.Exit(reasonInternalError)
			}
			fmt.Println("saving text success")
		} else {
			if err := app.SaveTextDataFromFile(ctx, textEdit.file, textEdit.name, textEdit.id); err != nil {
				fmt.Printf("saving text failed: %s\n", err)
				os.Exit(reasonInternalError)
			}
			fmt.Println("saving text success")
		}
	},
}

type textDeleteFlags struct {
	id int
}

var textDelete = textDeleteFlags{}

var textDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		if textDelete.id == 0 {
			fmt.Println("'id' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.DeleteTextData(ctx, textDelete.id); err != nil {
			fmt.Printf("deleting bank card failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("deleting text data is completed")
	},
}

func init() {
	textAddCmd.PersistentFlags().StringVarP(&textAdd.data, "data", "", "", "data")
	textAddCmd.PersistentFlags().StringVarP(&textAdd.file, "file", "", "", "file")
	textAddCmd.PersistentFlags().StringVarP(&textAdd.name, "name", "", "", "name")

	textGetCmd.PersistentFlags().IntVarP(&textGet.id, "id", "", 0, "id")

	textEditCmd.PersistentFlags().IntVarP(&textEdit.id, "id", "", 0, "id")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.data, "data", "", "", "data")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.name, "name", "", "", "name")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.file, "file", "", "", "file")

	textDeleteCmd.PersistentFlags().IntVarP(&textDelete.id, "id", "", 0, "id")

	textCmd.AddCommand(textAddCmd)
	textCmd.AddCommand(textGetCmd)
	textCmd.AddCommand(textEditCmd)
	textCmd.AddCommand(textDeleteCmd)
}
