package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bankCmd = &cobra.Command{
	Use:   "account",
	Short: "",
	Long:  ``,
}

type bankAddFlags struct {
	number  string
	expires string
	cvv     int
}

var bankAdd = bankAddFlags{}

var bankAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called bank add")
	},
}

type bankGetFlags struct {
	id int
}

var bankGet = bankGetFlags{}

var bankGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called bank get")
	},
}

type bankEditFlags struct {
	id      int
	number  string
	expires string
	cvv     int
}

var bankEdit = bankEditFlags{}

var bankEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called bank edit")
	},
}

type bankDeleteFlags struct {
	id int
}

var bankDelete = bankDeleteFlags{}

var bankDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called bank delete")
	},
}

func init() {
	bankAddCmd.PersistentFlags().StringVarP(&bankAdd.number, "number", "", "", "number")
	bankAddCmd.PersistentFlags().StringVarP(&bankAdd.expires, "expires", "", "", "date of expiration, format $month/$year, e.g 01/2020")
	bankAddCmd.PersistentFlags().IntVarP(&bankAdd.cvv, "cvv", "", 0, "cvv")

	bankGetCmd.PersistentFlags().IntVarP(&bankGet.id, "id", "", 0, "id")

	bankEditCmd.PersistentFlags().IntVarP(&bankEdit.id, "id", "", 0, "id")
	bankEditCmd.PersistentFlags().StringVarP(&bankEdit.number, "number", "", "", "number")
	bankEditCmd.PersistentFlags().StringVarP(&bankEdit.expires, "expires", "", "", "expires")
	bankEditCmd.PersistentFlags().IntVarP(&bankEdit.cvv, "cvv", "", 0, "cvv")

	bankDeleteCmd.PersistentFlags().IntVarP(&bankDelete.id, "id", "", 0, "id")

	bankCmd.AddCommand(bankAddCmd)
	bankCmd.AddCommand(bankGetCmd)
	bankCmd.AddCommand(bankEditCmd)
	bankCmd.AddCommand(bankDeleteCmd)
}
