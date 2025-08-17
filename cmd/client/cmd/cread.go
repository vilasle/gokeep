package cmd

import (
	"context"
	"fmt"
	"os"

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
		app := initCLIClient()
		defer app.Close()

		if creadAdd.login == "" {
			fmt.Println("login is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if creadAdd.password == "" {
			fmt.Println("password is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.SaveLoginPassword(ctx, creadAdd.login, creadAdd.password, 0); err != nil {
			fmt.Printf("saving login password failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("saving login password success")
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
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.GetLoginPassword(ctx, creadGet.id); err != nil {
			fmt.Printf("getting login failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
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
		app := initCLIClient()
		defer app.Close()

		if creadEdit.login == "" {
			fmt.Println("'login' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if creadEdit.password == "" {
			fmt.Println("'password' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		if creadEdit.id == 0 {
			fmt.Println("'id' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.SaveLoginPassword(ctx, creadEdit.login, creadEdit.password, creadEdit.id); err != nil {
			fmt.Printf("saving login password failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
		fmt.Println("saving login password success")
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
		app := initCLIClient()
		defer app.Close()

		if creadDelete.id == 0 {
			fmt.Println("'id' is required")
			os.Exit(reasonNotFillRequiredArgs)
		}

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.DeleteLoginPassword(ctx, creadDelete.id); err != nil {
			fmt.Printf("deleting login password failed: %s\n", err)
			os.Exit(reasonInternalError)
		}
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
