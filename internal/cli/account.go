package cli

import (
	"fmt"
	"strconv"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
)

func newCreateAccountCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "create-account <account_number>",
		Short: "Create a new account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account number: %w", err)
			}
			return tigerBeagle.CreateAccount(id)
		},
	}
}

func newGetAccountCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "get-account <account_number>",
		Short: "Get account details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account number: %w", err)
			}
			account, err := tigerBeagle.GetAccount(id)
			if err != nil {
				return err
			}
			fmt.Printf("Account details: %+v\n", account)
			return nil
		},
	}
}

func newMigrateAccountsCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-accounts <json_file>",
		Short: "Migrate accounts from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tigerBeagle.MigrateAccounts(args[0])
		},
	}
}
