package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/spf13/cobra"
)

func newDoctorCmd(tigerBeagle *app.TigerBeagle) *cobra.Command {
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

			fmt.Printf("Attempting to connect to TigerBeetle (timeout: %s)\n", totalTimeout)

			errChan := make(chan error, 1)
			go func() {
				errChan <- tigerBeagle.ValidateConnectivity()
			}()

			select {
			case err := <-errChan:
				if err == nil {
					fmt.Println("Successfully connected to TigerBeetle")
					return nil
				}
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, "client version is too old"):
					fmt.Println("Connection failed: Client version is too old")
					fmt.Println("Please update your TigerBeetle client")
					fmt.Printf("Error details: %v\n", err)
					return fmt.Errorf("connection failed due to outdated client version")
				case strings.Contains(errStr, "session was evicted"):
					fmt.Println("Connection failed: Session was evicted")
					fmt.Println("This might be due to a version mismatch or other issues")
					fmt.Printf("Error details: %v\n", err)
					return fmt.Errorf("connection failed due to session eviction")
				default:
					fmt.Printf("Failed to connect to TigerBeetle: %v\n", err)
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
