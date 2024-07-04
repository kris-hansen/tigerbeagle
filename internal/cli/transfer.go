package cli

import (
	"fmt"
	"strconv"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
)

func newTransferCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer <debit_account> <credit_account> <amount>",
		Short: "Transfer funds between accounts",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			debit, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid debit account: %w", err)
			}
			credit, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid credit account: %w", err)
			}
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			return tigerBeagle.Transfer(debit, credit, amount)
		},
	}
}

func newBulkTransferCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "bulk-transfer <debit_account> <credit_account> <amount> <iterations>",
		Short: "Perform multiple transfers in bulk",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			debit, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid debit account: %w", err)
			}
			credit, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid credit account: %w", err)
			}
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			iterations, err := strconv.Atoi(args[3])
			if err != nil {
				return fmt.Errorf("invalid number of iterations: %w", err)
			}
			return tigerBeagle.BulkTransfer(iterations, debit, credit, amount)
		},
	}
}

func newMigrateTransfersCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-transfers <json_file>",
		Short: "Migrate transfers from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tigerBeagle.MigrateTransfers(args[0])
		},
	}
}
