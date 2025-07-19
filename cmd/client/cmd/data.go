/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	dataCmd.PersistentFlags().StringArrayVarP(&meta, "meta", "m", []string{}, "meta data")

	dataCmd.AddCommand(creadCmd)
	dataCmd.AddCommand(bankCmd)
	dataCmd.AddCommand(textCmd)
	dataCmd.AddCommand(binaryCmd)

}
