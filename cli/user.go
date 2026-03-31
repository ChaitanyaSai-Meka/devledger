package cli

import (
	"encoding/json"
	"fmt"
	"net/url"

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
		body, err := json.Marshal(map[string]string{"username": username})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		resp, err := post("/users", string(body))
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

var listUserCmd = &cobra.Command{
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

var getUserGroupCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups a user belongs to",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := get("/users/" + url.PathEscape(username) + "/groups")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		resp, err := deleteReq("/users/" + url.PathEscape(username))
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
	userCmd.AddCommand(listUserCmd)
	userCmd.AddCommand(getUserGroupCmd)
	userCmd.AddCommand(deleteUserCmd)
	createUserCmd.Flags().StringP("username", "u", "", "username for the new user")
	if err := createUserCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	getUserGroupCmd.Flags().StringP("username", "u", "", "username to list groups for")
	if err := getUserGroupCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	deleteUserCmd.Flags().StringP("username", "u", "", "username to delete")
	if err := deleteUserCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
}
