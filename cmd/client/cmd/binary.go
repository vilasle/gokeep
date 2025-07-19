package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "",
	Long:  ``,
}

type binaryAddFlags struct {
	file string
}

var binaryAdd = binaryAddFlags{}

var binaryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called binary add")
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
		fmt.Println("called binary get")
	},
}

type binaryEditFlags struct {
	id   int
	file string
}

var binaryEdit = binaryEditFlags{}

var binaryEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called binary edit")
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
		fmt.Println("called binary delete")
	},
}

func init() {
	binaryAddCmd.PersistentFlags().StringVarP(&binaryAdd.file, "file", "", "", "file")

	binaryGetCmd.PersistentFlags().IntVarP(&binaryGet.id, "id", "", 0, "id")

	binaryEditCmd.PersistentFlags().IntVarP(&binaryEdit.id, "id", "", 0, "id")
	binaryEditCmd.PersistentFlags().StringVarP(&binaryEdit.file, "file", "", "", "file")

	binaryDeleteCmd.PersistentFlags().IntVarP(&binaryDelete.id, "id", "", 0, "id")

	binaryCmd.AddCommand(binaryAddCmd)
	binaryCmd.AddCommand(binaryGetCmd)
	binaryCmd.AddCommand(binaryEditCmd)
	binaryCmd.AddCommand(binaryDeleteCmd)
}