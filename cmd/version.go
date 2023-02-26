package main

import (
	"fmt"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/info"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version info",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("%-18s %-18s Commit:%s (%s)\n",
		info.Title,
		info.Version,
		info.Commit,
		info.BuildTime)
}
