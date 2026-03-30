package cli

import (
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "View balances",
}