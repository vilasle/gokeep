package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var creadCmd = &cobra.Command{
	Use:   "cread",
	Short: "",
	Long:  ``,
}

type credAddFlags struct {
	login    string
	password string
}

var creadAdd = credAddFlags{}

var creadAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called cread add")
	},
}

type createGetFlags struct {
	login string
	id    int
}

var creadGet = createGetFlags{}

var creadGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called cread get")
	},
}

type credEditFlags struct {
	id       int
	login    string
	password string
}

var creadEdit = credEditFlags{}

var creadEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called cread edit")
	},
}

type credDeleteFlags struct {
	id int
}

var creadDelete = credDeleteFlags{}

var creadDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("called cread delete")
	},
}

func init() {
	creadAddCmd.PersistentFlags().StringVarP(&creadAdd.login, "login", "", "", "login")
	creadAddCmd.PersistentFlags().StringVarP(&creadAdd.password, "password", "", "", "password")

	creadGetCmd.PersistentFlags().StringVarP(&creadGet.login, "login", "", "", "login")
	creadGetCmd.PersistentFlags().IntVarP(&creadGet.id, "id", "", 0, "id")

	creadEditCmd.PersistentFlags().IntVarP(&creadEdit.id, "id", "", 0, "login")
	creadEditCmd.PersistentFlags().StringVarP(&creadEdit.login, "login", "", "", "login")
	creadEditCmd.PersistentFlags().StringVarP(&creadEdit.password, "password", "", "", "password")

	creadDeleteCmd.PersistentFlags().IntVarP(&creadDelete.id, "id", "", 0, "id")

	creadCmd.AddCommand(creadAddCmd)
	creadCmd.AddCommand(creadGetCmd)
	creadCmd.AddCommand(creadEditCmd)
	creadCmd.AddCommand(creadDeleteCmd)
}
