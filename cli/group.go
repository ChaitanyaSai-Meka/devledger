package cli

import (
	"fmt"

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
		resp, err := post("/groups", fmt.Sprintf(`{"groupname":"%s"}`, groupname))
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
		resp, err := get("/groups/" + groupname + "/members")
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
		resp, err := deleteReq("/groups/" + groupname + "/members/" + username)
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
		resp, err := deleteReq("/groups/" + groupname)
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
		resp, err := post("/groups/"+groupname+"/members", fmt.Sprintf(`{"username":"%s"}`, username))
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
	createGroupCmd.MarkFlagRequired("groupname")
	getGroupMembersCmd.Flags().StringP("groupname", "g", "", "group name")
	getGroupMembersCmd.MarkFlagRequired("groupname")
	removeMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	removeMemberCmd.Flags().StringP("username", "u", "", "username to remove")
	removeMemberCmd.MarkFlagRequired("groupname")
	removeMemberCmd.MarkFlagRequired("username")
	deleteGroupCmd.Flags().StringP("groupname", "g", "", "group name to delete")
	deleteGroupCmd.MarkFlagRequired("groupname")
	addMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	addMemberCmd.Flags().StringP("username", "u", "", "username to add")
	addMemberCmd.MarkFlagRequired("groupname")
	addMemberCmd.MarkFlagRequired("username")
}
