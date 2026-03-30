package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := post("/users", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("User created successfully")
		} else {
			printResponse(resp)
		}
	},
}

var ListUserCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := get("/users")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var GetUserGroupCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups a user belongs to",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := get("/users/" + username + "/groups")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var DeleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := deleteReq("/users/" + username)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Printf("User '%s' deleted successfully\n", username)
		} else {
			printResponse(resp)
		}
	},
}

func init() {
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(ListUserCmd)
	userCmd.AddCommand(GetUserGroupCmd)
	userCmd.AddCommand(DeleteUserCmd)
	createUserCmd.Flags().StringP("username", "u", "", "username for the new user")
	createUserCmd.MarkFlagRequired("username")
	GetUserGroupCmd.Flags().StringP("username", "u", "", "username to list groups for")
	GetUserGroupCmd.MarkFlagRequired("username")
	DeleteUserCmd.Flags().StringP("username", "u", "", "username to delete")
	DeleteUserCmd.MarkFlagRequired("username")
}
