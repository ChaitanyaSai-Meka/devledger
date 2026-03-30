package cli

import (
	"github.com/spf13/cobra"
)

var expenseCmd = &cobra.Command{
	Use: "expense",
	Short: "Manage expenses",
}