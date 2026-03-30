package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var expenseCmd = &cobra.Command{
	Use:   "expense",
	Short: "Manage expenses",
}

var addExpenseCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an expense to a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		description, _ := cmd.Flags().GetString("description")
		amount, _ := cmd.Flags().GetString("amount")
		paidby, _ := cmd.Flags().GetString("paidby")
		resp, err := post("/groups/"+groupname+"/expenses", fmt.Sprintf(`{"description":"%s","amount":"%s","paid_by":"%s"}`, description, amount, paidby))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("Expense added successfully")
		} else {
			printResponse(resp)
		}
	},
}

var listExpensebyGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "List expenses in a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		resp, err := get("/groups/" + groupname + "/expenses")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var listExpenseByUserCmd = &cobra.Command{
	Use:   "listbyuser",
	Short: "List expenses paid by a user",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := get("/users/" + username + "/expenses")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var deleteExpenseCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an expense by ID",
	Run: func(cmd *cobra.Command, args []string) {
		expenseID, _ := cmd.Flags().GetInt("expenseid")
		resp, err := deleteReq("/expenses/" + strconv.Itoa(expenseID))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Println("Expense deleted successfully")
		} else {
			printResponse(resp)
		}
	},
}

var settleExpenseCmd = &cobra.Command{
	Use:   "settle",
	Short: "Settle an expense for a user",
	Run: func(cmd *cobra.Command, args []string) {
		expenseID, _ := cmd.Flags().GetInt("expenseid")
		username, _ := cmd.Flags().GetString("username")
		resp, err := post("/expenses/"+strconv.Itoa(expenseID)+"/settle/"+username, "")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Printf("Expense %d settled for user '%s'\n", expenseID, username)
		} else {
			printResponse(resp)
		}
	},
}

var listUnsettledSplitsCmd = &cobra.Command{
	Use:   "unsettled",
	Short: "List unsettled splits for a user",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := get("/users/" + username + "/unsettled-splits")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var expenseInDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Get expense details by ID",
	Run: func(cmd *cobra.Command, args []string) {
		expenseID, _ := cmd.Flags().GetInt("expenseid")
		resp, err := get("/expenses/" + strconv.Itoa(expenseID) + "/detail")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

func init() {
	expenseCmd.AddCommand(addExpenseCmd)
	expenseCmd.AddCommand(listExpensebyGroupCmd)
	expenseCmd.AddCommand(listExpenseByUserCmd)
	expenseCmd.AddCommand(deleteExpenseCmd)
	expenseCmd.AddCommand(settleExpenseCmd)
	expenseCmd.AddCommand(listUnsettledSplitsCmd)
	expenseCmd.AddCommand(expenseInDetailCmd)

	addExpenseCmd.Flags().StringP("groupname", "g", "", "group name")
	addExpenseCmd.Flags().StringP("description", "d", "", "expense description")
	addExpenseCmd.Flags().StringP("amount", "a", "", "expense amount")
	addExpenseCmd.Flags().StringP("paidby", "p", "", "username who paid")
	addExpenseCmd.MarkFlagRequired("groupname")
	addExpenseCmd.MarkFlagRequired("description")
	addExpenseCmd.MarkFlagRequired("amount")
	addExpenseCmd.MarkFlagRequired("paidby")

	listExpensebyGroupCmd.Flags().StringP("groupname", "g", "", "group name")
	listExpensebyGroupCmd.MarkFlagRequired("groupname")

	listExpenseByUserCmd.Flags().StringP("username", "u", "", "username")
	listExpenseByUserCmd.MarkFlagRequired("username")

	deleteExpenseCmd.Flags().IntP("expenseid", "e", 0, "expense ID")
	deleteExpenseCmd.MarkFlagRequired("expenseid")

	settleExpenseCmd.Flags().IntP("expenseid", "e", 0, "expense ID")
	settleExpenseCmd.Flags().StringP("username", "u", "", "username")
	settleExpenseCmd.MarkFlagRequired("expenseid")
	settleExpenseCmd.MarkFlagRequired("username")

	listUnsettledSplitsCmd.Flags().StringP("username", "u", "", "username")
	listUnsettledSplitsCmd.MarkFlagRequired("username")

	expenseInDetailCmd.Flags().IntP("expenseid", "e", 0, "expense ID")
	expenseInDetailCmd.MarkFlagRequired("expenseid")
}
