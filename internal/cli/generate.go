package cli

import (
	"fmt"
	"strconv"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGenerateCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [account|transfer] <number>",
		Short: "Generate sample JSON files for accounts or transfers",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			generateType := args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid number: %w", err)
			}

			ledger := viper.GetUint32("ledger")
			code := uint16(viper.GetUint32("code"))
			flags := uint16(viper.GetUint32("flags"))

			switch generateType {
			case "account":
				return tigerBeagle.GenerateAccounts(number, ledger, code, flags)
			case "transfer":
				return tigerBeagle.GenerateTransfers(number, ledger, code, flags)
			default:
				return fmt.Errorf("invalid generate type: must be 'account' or 'transfer'")
			}
		},
	}

	return cmd
}
