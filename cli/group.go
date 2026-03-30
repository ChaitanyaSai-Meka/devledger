package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
}

var CreateGroupCmd = &cobra.Command{
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

var ListGroupCmd = &cobra.Command{
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

var GetGroupMembersCmd = &cobra.Command{
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

var RemoveMemberCmd = &cobra.Command{
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

var DeleteGroupCmd = &cobra.Command{
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

var AddMemberCmd = &cobra.Command{
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
	groupCmd.AddCommand(CreateGroupCmd)
	groupCmd.AddCommand(ListGroupCmd)
	groupCmd.AddCommand(GetGroupMembersCmd)
	groupCmd.AddCommand(RemoveMemberCmd)
	groupCmd.AddCommand(DeleteGroupCmd)
	groupCmd.AddCommand(AddMemberCmd)

	CreateGroupCmd.Flags().StringP("groupname", "g", "", "name for the new group")
	CreateGroupCmd.MarkFlagRequired("groupname")
	GetGroupMembersCmd.Flags().StringP("groupname", "g", "", "group name")
	GetGroupMembersCmd.MarkFlagRequired("groupname")
	RemoveMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	RemoveMemberCmd.Flags().StringP("username", "u", "", "username to remove")
	RemoveMemberCmd.MarkFlagRequired("groupname")
	RemoveMemberCmd.MarkFlagRequired("username")
	DeleteGroupCmd.Flags().StringP("groupname", "g", "", "group name to delete")
	DeleteGroupCmd.MarkFlagRequired("groupname")
	AddMemberCmd.Flags().StringP("groupname", "g", "", "group name")
	AddMemberCmd.Flags().StringP("username", "u", "", "username to add")
	AddMemberCmd.MarkFlagRequired("groupname")
	AddMemberCmd.MarkFlagRequired("username")
}
