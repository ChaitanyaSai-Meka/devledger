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
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		body, err := json.Marshal(map[string]string{"groupname": groupname})
		if err != nil {
			return err
		}
		resp, err := post("/groups", string(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("Group created successfully")
			return nil
		}
		return printResponse(resp)
	},
}

var listGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := get("/groups")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return printResponse(resp)
	},
}

var getGroupMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "List members of a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		resp, err := get("/groups/" + url.PathEscape(groupname) + "/members")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return printResponse(resp)
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		resp, err := deleteReq("/groups/" + url.PathEscape(groupname) + "/members/" + url.PathEscape(username))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return printResponse(resp)
	},
}

var deleteGroupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		resp, err := deleteReq("/groups/" + url.PathEscape(groupname))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Printf("Group '%s' deleted successfully\n", groupname)
			return nil
		}
		return printResponse(resp)
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupname, err := cmd.Flags().GetString("groupname")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		body, err := json.Marshal(map[string]string{"username": username})
		if err != nil {
			return err
		}
		resp, err := post("/groups/"+url.PathEscape(groupname)+"/members", string(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 201 {
			fmt.Println("Member added successfully")
			return nil
		}
		return printResponse(resp)
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
