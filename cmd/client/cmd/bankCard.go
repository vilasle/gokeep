package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vilasle/gokeep/internal/client"
)

var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "manager of bank cards",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type bankAddFlags struct {
	number  string
	expires string
	cvv     int
}

var bankAdd = bankAddFlags{}

var bankAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add bank card on server and on local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := bankCardAddHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

type bankGetFlags struct {
	id int
}

var bankGet = bankGetFlags{}

var bankGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get bank card from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := bankCardGetHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
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
	Short: "edit bank card on local storage and on server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := bankCardEditHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
	},
}

type bankDeleteFlags struct {
	id int
}

var bankDelete = bankDeleteFlags{}

var bankDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete bank card on local storage and on server",
	Run: func(cmd *cobra.Command, args []string) {
		app := initCLIClient()
		defer app.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if code := bankCardDeleteHandle(ctx, app); code != 0 {
			cmd.Usage()
			os.Exit(code)
		}
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

func bankCardAddHandle(ctx context.Context, app client.Client) int {
	if bankAdd.number == "" {
		fmt.Println("--number argument is required")
		return reasonNotFillRequiredArgs
	}

	if bankAdd.cvv == 0 {
		fmt.Println("--cvv argument is required")
		return reasonNotFillRequiredArgs
	}

	if bankAdd.expires == "" {
		fmt.Println("--expires argument is required")
		return reasonNotFillRequiredArgs
	}

	metadata := prepareMetadata()

	if err := app.SaveBankCard(ctx, bankAdd.number, bankAdd.expires, bankAdd.cvv, 0, metadata); err != nil {
		fmt.Printf("saving bank card failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("saving bank card success")
	return 0
}

func bankCardEditHandle(ctx context.Context, app client.Client) int {
	if bankEdit.number == "" {
		fmt.Println("--number argument is required")
		return reasonNotFillRequiredArgs
	}

	if bankEdit.cvv == 0 {
		fmt.Println("--cvv argument is required")
		return reasonNotFillRequiredArgs
	}

	if bankEdit.expires == "" {
		fmt.Println("--expires argument is required")
		return reasonNotFillRequiredArgs
	}

	if bankEdit.id == 0 {
		fmt.Println("--id argument is required")
		return reasonNotFillRequiredArgs
	}

	metadata := prepareMetadata()

	if err := app.SaveBankCard(ctx, bankEdit.number, bankEdit.expires, bankEdit.cvv, bankEdit.id, metadata); err != nil {
		fmt.Printf("saving bank card failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("saving bank card success")
	return 0
}

func bankCardGetHandle(ctx context.Context, app client.Client) int {
	if err := app.GetBankCard(ctx, bankGet.id); err != nil {
		fmt.Printf("getting bank card failed: %s\n", err)
		return reasonInternalError
	}
	return 0
}

func bankCardDeleteHandle(ctx context.Context, app client.Client) int {
	if bankDelete.id == 0 {
		fmt.Println("--id argument is required")
		return reasonNotFillRequiredArgs
	}

	if err := app.DeleteBankCard(ctx, bankDelete.id); err != nil {
		fmt.Printf("deleting bank card failed: %s\n", err)
		return reasonInternalError
	}
	fmt.Println("deleting band card is completed")
	return 0
}
