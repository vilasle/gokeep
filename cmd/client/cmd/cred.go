package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
)

var credCmd = &cobra.Command{
	Use:   "cred",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type credAddFlags struct {
	login    string
	password string
}

var credAdd = credAddFlags{}

var credAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		credAddHandle(ctx, app)
	},
}

type createGetFlags struct {
	login string
	id    int
}

var credGet = createGetFlags{}

var credGetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		//TODO add waiting SIGNAL and cancel if got it
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		credGetHandle(ctx, app)
	},
}

type credEditFlags struct {
	id       int
	login    string
	password string
}

var credEdit = credEditFlags{}

var credEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		credEditHandle(ctx, app)
	},
}

type credDeleteFlags struct {
	id int
}

var credDelete = credDeleteFlags{}

var credDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		credDeleteHandle(ctx, app)
	},
}

func init() {
	credAddCmd.PersistentFlags().StringVarP(&credAdd.login, "login", "", "", "login")
	credAddCmd.PersistentFlags().StringVarP(&credAdd.password, "password", "", "", "password")

	credGetCmd.PersistentFlags().StringVarP(&credGet.login, "login", "", "", "login")
	credGetCmd.PersistentFlags().IntVarP(&credGet.id, "id", "", 0, "id")

	credEditCmd.PersistentFlags().IntVarP(&credEdit.id, "id", "", 0, "login")
	credEditCmd.PersistentFlags().StringVarP(&credEdit.login, "login", "", "", "login")
	credEditCmd.PersistentFlags().StringVarP(&credEdit.password, "password", "", "", "password")

	credDeleteCmd.PersistentFlags().IntVarP(&credDelete.id, "id", "", 0, "id")

	credCmd.AddCommand(credAddCmd)
	credCmd.AddCommand(credGetCmd)
	credCmd.AddCommand(credEditCmd)
	credCmd.AddCommand(credDeleteCmd)
}

func credAddHandle(ctx context.Context, app client.Client) {
	if credAdd.login == "" {
		fmt.Println("--login argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}

	if credAdd.password == "" {
		fmt.Println("--password argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}
	metadata := prepareMetadata()

	if err := app.SaveLoginPassword(ctx, credAdd.login, credAdd.password, 0, metadata); err != nil {
		fmt.Printf("saving login password failed: %s\n", err)
		os.Exit(reasonInternalError)
	}
	fmt.Println("saving login password success")
}

func credEditHandle(ctx context.Context, app client.Client) {
	if credEdit.login == "" {
		fmt.Println("--login argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}

	if credEdit.password == "" {
		fmt.Println("--password argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}

	if credEdit.id == 0 {
		fmt.Println("--id argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}

	metadata := prepareMetadata()

	if err := app.SaveLoginPassword(ctx, credEdit.login, credEdit.password, credEdit.id, metadata); err != nil {
		fmt.Printf("saving login password failed: %s\n", err)
		os.Exit(reasonInternalError)
	}
	fmt.Println("saving login password success")
}

func credGetHandle(ctx context.Context, app client.Client) {
	if err := app.GetLoginPassword(ctx, credGet.id); err != nil {
		fmt.Printf("getting login failed: %s\n", err)
		os.Exit(reasonInternalError)
	}
}

func credDeleteHandle(ctx context.Context, app client.Client) {
	if credDelete.id == 0 {
		fmt.Println("--id argument is required")
		os.Exit(reasonNotFillRequiredArgs)
	}

	if err := app.DeleteLoginPassword(ctx, credDelete.id); err != nil {
		fmt.Printf("deleting login password failed: %s\n", err)
		os.Exit(reasonInternalError)
	}
	fmt.Println("deleting login password is completed")
}
