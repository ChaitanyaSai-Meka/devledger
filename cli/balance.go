package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "View balances",
}

var groupBalanceCmd = &cobra.Command{
	Use:   "group",
	Short: "View group balances",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		resp, err := get("/groups/" + groupname + "/balance")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var simplifyGroupBalanceCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify group balances",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		resp, err := get("/groups/" + groupname + "/balance/simplify")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

func init() {
	balanceCmd.AddCommand(groupBalanceCmd)
	balanceCmd.AddCommand(simplifyGroupBalanceCmd)
	groupBalanceCmd.Flags().StringP("groupname", "g", "", "Name of the group")
	groupBalanceCmd.MarkFlagRequired("groupname")
	simplifyGroupBalanceCmd.Flags().StringP("groupname", "g", "", "Name of the group")
	simplifyGroupBalanceCmd.MarkFlagRequired("groupname")
}
