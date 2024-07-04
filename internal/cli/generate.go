package cli

import (
	"fmt"
	"strconv"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
)

func newGenerateCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [account|transaction] <number>",
		Short: "Generate sample JSON files for accounts or transactions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			generateType := args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid number: %w", err)
			}

			switch generateType {
			case "account":
				return tigerBeagle.GenerateAccounts(number)
			case "transaction":
				return tigerBeagle.GenerateTransfers(number)
			default:
				return fmt.Errorf("invalid generate type: must be 'account' or 'transaction'")
			}
		},
	}

	return cmd
}
