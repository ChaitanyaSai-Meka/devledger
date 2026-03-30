package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devledger",
	Short: "A developer expense splitting tool",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(expenseCmd)
	rootCmd.AddCommand(balanceCmd)
}
