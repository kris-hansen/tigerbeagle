package cli

import (
	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand(tigerBeagle *app.TigerBeagle) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tigerbeagle",
		Short: "TigerBeagle is a CLI tool for TigerBeetle ledger data management",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return tigerBeagle.InitClient(viper.GetString("tb_address"))
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			tigerBeagle.CloseClient()
		},
	}

	rootCmd.PersistentFlags().String("tb-address", "3000", "TigerBeetle address")
	viper.BindPFlag("tb_address", rootCmd.PersistentFlags().Lookup("tb-address"))

	// Account commands
	rootCmd.AddCommand(
		newCreateAccountCmd(tigerBeagle),
		newGetAccountCmd(tigerBeagle),
		newMigrateAccountsCmd(tigerBeagle),
	)

	// Transfer commands
	rootCmd.AddCommand(
		newTransferCmd(tigerBeagle),
		newBulkTransferCmd(tigerBeagle),
		newMigrateTransfersCmd(tigerBeagle),
	)

	// Other commands
	rootCmd.AddCommand(
		newDoctorCmd(tigerBeagle),
		newGenerateCmd(tigerBeagle),
	)

	return rootCmd
}
