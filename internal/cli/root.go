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
	rootCmd.PersistentFlags().Uint32("ledger", 700, "Ledger ID")
	rootCmd.PersistentFlags().Uint16("code", 10, "Account/Transfer code")
	rootCmd.PersistentFlags().Uint16("flags", 0, "Account/Transfer flags")

	viper.BindPFlag("tb_address", rootCmd.PersistentFlags().Lookup("tb-address"))
	viper.BindPFlag("ledger", rootCmd.PersistentFlags().Lookup("ledger"))
	viper.BindPFlag("code", rootCmd.PersistentFlags().Lookup("code"))
	viper.BindPFlag("flags", rootCmd.PersistentFlags().Lookup("flags"))

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
