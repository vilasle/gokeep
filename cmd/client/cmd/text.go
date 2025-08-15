package cmd

import (
	"fmt"

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
		fmt.Println("called text add")
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
		fmt.Println("called text get")
	},
}

type textEditFlags struct {
	id   int
	data string
	file string
}

var textEdit = textEditFlags{}

var textEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called text edit")
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
		fmt.Println("called text delete")
	},
}

func init() {
	textAddCmd.PersistentFlags().StringVarP(&textAdd.data, "data", "", "", "data")
	textAddCmd.PersistentFlags().StringVarP(&textAdd.file, "file", "", "", "file")
	textAddCmd.PersistentFlags().StringVarP(&textAdd.name, "name", "", "", "name")

	textGetCmd.PersistentFlags().IntVarP(&textGet.id, "id", "", 0, "id")

	textEditCmd.PersistentFlags().IntVarP(&textEdit.id, "id", "", 0, "id")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.data, "data", "", "", "data")
	textEditCmd.PersistentFlags().StringVarP(&textEdit.file, "file", "", "", "file")

	textDeleteCmd.PersistentFlags().IntVarP(&textDelete.id, "id", "", 0, "id")

	textCmd.AddCommand(textAddCmd)
	textCmd.AddCommand(textGetCmd)
	textCmd.AddCommand(textEditCmd)
	textCmd.AddCommand(textDeleteCmd)
}
