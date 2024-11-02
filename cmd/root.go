package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "jondev-server",
	Short: "Backend for jondev, a simple and modern application to build dev & designer portfolio",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
