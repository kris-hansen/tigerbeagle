package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
)

func newDoctorCmd(tigerBeagle app.TigerBeagleInterface) *cobra.Command {
	var attempts int
	var timeoutPerAttempt int

	cmd := &cobra.Command{
		Use:          "doctor",
		Short:        "Validate the connectivity to TigerBeetle",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			totalTimeout := time.Duration(attempts*timeoutPerAttempt) * time.Second
			ctx, cancel := context.WithTimeout(cmd.Context(), totalTimeout)
			defer cancel()

			fmt.Fprintf(cmd.OutOrStdout(), "Attempting to connect to TigerBeetle (timeout: %s)\n", totalTimeout)

			errChan := make(chan error, 1)
			go func() {
				errChan <- tigerBeagle.ValidateConnectivity()
			}()

			select {
			case err := <-errChan:
				if err == nil {
					fmt.Fprintln(cmd.OutOrStdout(), "Successfully connected to TigerBeetle")
					return nil
				}
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, "client version is too old"):
					fmt.Fprintln(cmd.OutOrStdout(), "Connection failed: Client version is too old")
					fmt.Fprintln(cmd.OutOrStdout(), "Please update your TigerBeetle client")
					fmt.Fprintf(cmd.OutOrStdout(), "Error details: %v\n", err)
					return fmt.Errorf("connection failed due to outdated client version")
				case strings.Contains(errStr, "session was evicted"):
					fmt.Fprintln(cmd.OutOrStdout(), "Connection failed: Session was evicted")
					fmt.Fprintln(cmd.OutOrStdout(), "This might be due to a version mismatch or other issues")
					fmt.Fprintf(cmd.OutOrStdout(), "Error details: %v\n", err)
					return fmt.Errorf("connection failed due to session eviction")
				default:
					fmt.Fprintf(cmd.OutOrStdout(), "Failed to connect to TigerBeetle: %v\n", err)
					return fmt.Errorf("connection failed")
				}
			case <-ctx.Done():
				return fmt.Errorf("operation timed out after %s", totalTimeout)
			}
		},
	}

	cmd.Flags().IntVarP(&attempts, "attempts", "a", 5, "Number of logical attempts (used to calculate total timeout)")
	cmd.Flags().IntVarP(&timeoutPerAttempt, "timeout-per-attempt", "t", 10, "Timeout in seconds per attempt")

	return cmd
}
