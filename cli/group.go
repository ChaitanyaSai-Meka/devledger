package cli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
}

var createGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		body, err := json.Marshal(map[string]string{"groupname": groupname})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		resp, err := post("/groups", string(body))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("Group created successfully")
		} else {
			printResponse(resp)
		}
	},
}

var listGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := get("/groups")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var getGroupMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "List members of a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		resp, err := get("/groups/" + url.PathEscape(groupname) + "/members")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		username, _ := cmd.Flags().GetString("username")
		resp, err := deleteReq("/groups/" + url.PathEscape(groupname) + "/members/" + url.PathEscape(username))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		printResponse(resp)
	},
}

var deleteGroupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		resp, err := deleteReq("/groups/" + url.PathEscape(groupname))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Printf("Group '%s' deleted successfully\n", groupname)
		} else {
			printResponse(resp)
		}
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to a group",
	Run: func(cmd *cobra.Command, args []string) {
		groupname, _ := cmd.Flags().GetString("groupname")
		username, _ := cmd.Flags().GetString("username")
		body, err := json.Marshal(map[string]string{"username": username})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		resp, err := post("/groups/"+url.PathEscape(groupname)+"/members", string(body))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("Member added successfully")
		} else {
			printResponse(resp)
		}
	},
}

func init() {
	groupCmd.AddCommand(createGroupCmd)
	groupCmd.AddCommand(listGroupCmd)
	groupCmd.AddCommand(getGroupMembersCmd)
	groupCmd.AddCommand(removeMemberCmd)
	groupCmd.AddCommand(deleteGroupCmd)
	groupCmd.AddCommand(addMemberCmd)

	createGroupCmd.Flags().StringP("groupname", "g", "", "name for the new group")
	if err := createGroupCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	getGroupMembersCmd.Flags().StringP("groupname", "g", "", "group name")
	if err := getGroupMembersCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	removeMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	removeMemberCmd.Flags().StringP("username", "u", "", "username to remove")
	if err := removeMemberCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	if err := removeMemberCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	deleteGroupCmd.Flags().StringP("groupname", "g", "", "group name to delete")
	if err := deleteGroupCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	addMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	addMemberCmd.Flags().StringP("username", "u", "", "username to add")
	if err := addMemberCmd.MarkFlagRequired("groupname"); err != nil {
		panic(err)
	}
	if err := addMemberCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
}
