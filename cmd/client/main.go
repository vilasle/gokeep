/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cmp"

	"github.com/vilasle/gokeep/cmd/client/cmd"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	cmd.Version = cmp.Or(buildVersion, "N/A")
	cmd.Date = cmp.Or(buildDate, "N/A")
	cmd.Commit = cmp.Or(buildCommit, "N/A")

	cmd.Execute()
}
