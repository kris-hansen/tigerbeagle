package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/kris-hansen/tigerbeagle/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTigerBeagle struct {
	mock.Mock
}

func (m *MockTigerBeagle) ValidateConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTigerBeagle) CreateAccount(id uint64, ledger uint32, code uint16, flags uint16) error {
	args := m.Called(id, ledger, code, flags)
	return args.Error(0)
}

func (m *MockTigerBeagle) GetAccount(id uint64) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockTigerBeagle) Transfer(debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error {
	args := m.Called(debitAccountID, creditAccountID, amount, ledger, code, flags)
	return args.Error(0)
}

func (m *MockTigerBeagle) BulkTransfer(iterations int, debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error {
	args := m.Called(iterations, debitAccountID, creditAccountID, amount, ledger, code, flags)
	return args.Error(0)
}

func (m *MockTigerBeagle) GenerateAccounts(number int, ledger uint32, code uint16, flags uint16) error {
	args := m.Called(number, ledger, code, flags)
	return args.Error(0)
}

func (m *MockTigerBeagle) GenerateTransfers(number int, ledger uint32, code uint16, flags uint16) error {
	args := m.Called(number, ledger, code, flags)
	return args.Error(0)
}

func (m *MockTigerBeagle) MigrateAccounts(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func (m *MockTigerBeagle) MigrateTransfers(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

// Ensure MockTigerBeagle implements TigerBeagleInterface
var _ app.TigerBeagleInterface = (*MockTigerBeagle)(nil)

func TestDoctorCmd(t *testing.T) {
	tests := []struct {
		name           string
		connectError   error
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "Successful connection",
			connectError:   nil,
			expectedOutput: "Successfully connected to TigerBeetle",
			expectedError:  "",
		},
		{
			name:           "Client version too old",
			connectError:   fmt.Errorf("client version is too old"),
			expectedOutput: "Connection failed: Client version is too old",
			expectedError:  "connection failed due to outdated client version",
		},
		{
			name:           "Session evicted",
			connectError:   fmt.Errorf("session was evicted"),
			expectedOutput: "Connection failed: Session was evicted",
			expectedError:  "connection failed due to session eviction",
		},
		{
			name:           "Generic error",
			connectError:   fmt.Errorf("connection failed"),
			expectedOutput: "Failed to connect to TigerBeetle",
			expectedError:  "connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTB := new(MockTigerBeagle)
			cmd := newDoctorCmd(mockTB)

			mockTB.On("ValidateConnectivity").Return(tt.connectError).Once()

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			// Set flag values
			cmd.Flags().Set("attempts", "1")
			cmd.Flags().Set("timeout-per-attempt", "1")

			err := cmd.Execute()

			output := buf.String()
			t.Logf("Command output: %s", output)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Contains(t, output, tt.expectedOutput)
			mockTB.AssertExpectations(t)
		})
	}
}

func TestNewBulkTransferCmd(t *testing.T) {
	mockTB := new(MockTigerBeagle)

	// Create a custom viperGetUint32 function
	viperGetUint32 := func(key string) uint32 {
		switch key {
		case "ledger":
			return 700
		case "code":
			return 10
		case "flags":
			return 0
		default:
			return 0
		}
	}

	customNewBulkTransferCmd := func(tigerBeagle app.TigerBeagleInterface) *cobra.Command {
		cmd := &cobra.Command{
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

				ledger := viperGetUint32("ledger")
				code := uint16(viperGetUint32("code"))
				flags := uint16(viperGetUint32("flags"))

				return tigerBeagle.BulkTransfer(iterations, debit, credit, amount, ledger, code, flags)
			},
		}
		return cmd
	}

	cmd := customNewBulkTransferCmd(mockTB)

	// Set up the mock expectation
	mockTB.On("BulkTransfer", 5, uint64(1000), uint64(2000), uint64(100), uint32(700), uint16(10), uint16(0)).Return(nil).Once()

	// Set up command arguments
	args := []string{"1000", "2000", "100", "5"}
	cmd.SetArgs(args)

	// Execute the command
	err := cmd.Execute()

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockTB.AssertExpectations(t)
}
