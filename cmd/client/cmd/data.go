/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "",
	Long:  ``,
}

var meta []string

func init() {
	dataCmd.PersistentFlags().StringArrayVarP(&meta, "metadata", "m", []string{}, "metadata")

	dataCmd.AddCommand(credCmd)
	dataCmd.AddCommand(bankCmd)
	dataCmd.AddCommand(textCmd)
	dataCmd.AddCommand(binaryCmd)
	dataCmd.AddCommand(syncCmd)

}

func prepareMetadata() map[string]string {
	metadata := make(map[string]string)
	for _, v := range meta {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			fmt.Printf("invalid metadata: %s\n", v)
			fmt.Println("metadata must be in format key=value")
			os.Exit(reasonNotFillRequiredArgs)
		}
		metadata[kv[0]] = kv[1]
	}
	return metadata
}
