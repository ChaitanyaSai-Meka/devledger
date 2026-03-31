package cli

import (
	"fmt"
	"net/url"

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
		escapedgroupname := url.PathEscape(groupname)
		resp, err := get("/groups/" + escapedgroupname + "/balance")
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
		resp, err := get("/groups/" + url.PathEscape(groupname) + "/balance/simplify")
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
	if err := groupBalanceCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	simplifyGroupBalanceCmd.Flags().StringP("groupname", "g", "", "Name of the group")
	if err := simplifyGroupBalanceCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
}
