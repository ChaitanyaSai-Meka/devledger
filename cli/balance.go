package cli

import (
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
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		escapedgroupname := url.PathEscape(groupname)
		resp, err := get("/groups/" + escapedgroupname + "/balance")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return printResponse(resp)
	},
}

var simplifyGroupBalanceCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify group balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		resp, err := get("/groups/" + url.PathEscape(groupname) + "/balance/simplify")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return printResponse(resp)
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
