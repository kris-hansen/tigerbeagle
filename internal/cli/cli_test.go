package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/kris-hansen/tigerbeagle/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTigerBeagle mocks the TigerBeagleInterface
type MockTigerBeagle struct {
	mock.Mock
}

func (m *MockTigerBeagle) ValidateConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTigerBeagle) CreateAccount(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTigerBeagle) GetAccount(id uint64) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockTigerBeagle) Transfer(debitAccountID, creditAccountID, amount uint64) error {
	args := m.Called(debitAccountID, creditAccountID, amount)
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
