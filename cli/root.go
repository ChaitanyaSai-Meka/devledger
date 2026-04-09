package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "devledger",
	Short:        "A developer expense splitting tool",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprint(cmd.OutOrStdout(), renderLandingScreen())
		return err
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(expenseCmd)
	rootCmd.AddCommand(balanceCmd)
}

func renderLandingScreen() string {
	const (
		reset = "\033[0m"
		bold  = "\033[1m"
		dim   = "\033[2m"

		teal  = "\033[38;5;44m"
		mint  = "\033[38;5;86m"
		gold  = "\033[38;5;221m"
		coral = "\033[38;5;209m"
		slate = "\033[38;5;245m"
		sky   = "\033[38;5;117m"
	)

	var b strings.Builder

	writeRule := func(color string) {
		b.WriteString(color)
		b.WriteString(strings.Repeat("=", 72))
		b.WriteString(reset)
		b.WriteString("\n")
	}

	writeSection := func(title string, color string, lines []string) {
		b.WriteString(color)
		b.WriteString(bold)
		b.WriteString(title)
		b.WriteString(reset)
		b.WriteString("\n")
		for _, line := range lines {
			b.WriteString("  ")
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	writeRule(teal)
	b.WriteString(teal)
	b.WriteString(bold)
	b.WriteString("DEVLEDGER")
	b.WriteString(reset)
	b.WriteString("  ")
	b.WriteString(slate)
	b.WriteString("Local-first expense tracking for developer teams")
	b.WriteString(reset)
	b.WriteString("\n")
	writeRule(teal)
	b.WriteString("\n")

	writeSection("Command Areas", mint, []string{
		fmt.Sprintf("%suser%s     Create, list, and delete users", bold, reset),
		fmt.Sprintf("%sgroup%s    Create groups and manage members", bold, reset),
		fmt.Sprintf("%sexpense%s  Add expenses, inspect details, settle splits", bold, reset),
		fmt.Sprintf("%sbalance%s  View balances and suggested settlements", bold, reset),
	})

	writeSection("Quick Start", sky, []string{
		`devledger user create --username "alice"`,
		`devledger group create --groupname "backend-team"`,
		`devledger group add-member --groupname "backend-team" --username "alice"`,
		`devledger expense add --groupname "backend-team" --description "AWS Bill" --amount "30.00" --paidby "alice"`,
		`devledger balance simplify --groupname "backend-team"`,
	})

	writeSection("Useful Commands", gold, []string{
		`devledger user list`,
		`devledger group members --groupname "backend-team"`,
		`devledger expense list --groupname "backend-team"`,
		`devledger expense list-by-user --username "alice"`,
		`devledger expense unsettled --username "alice"`,
	})

	writeSection("Need More?", coral, []string{
		fmt.Sprintf("Run %sdevledger <command> --help%s for command-specific options", bold, reset),
		fmt.Sprintf("Example: %sdevledger expense --help%s", bold, reset),
	})

	b.WriteString(dim)
	b.WriteString("Server: http://localhost:38080")
	b.WriteString(reset)
	b.WriteString("\n")

	return b.String()
}
