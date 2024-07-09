package app

import (
	"fmt"
	"testing"

	"github.com/kris-hansen/tigerbeagle/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock tigerbeetle.Client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) CreateAccounts(accounts []models.Account) error {
	args := m.Called(accounts)
	return args.Error(0)
}

func (m *MockClient) LookupAccount(id uint64) (*models.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockClient) CreateTransfers(transfers []models.Transfer) error {
	args := m.Called(transfers)
	return args.Error(0)
}

func (m *MockClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClient) Close() {
	m.Called()
}

func TestValidateConnectivity(t *testing.T) {
	mockClient := new(MockClient)
	tb := &TigerBeagle{client: mockClient}

	tests := []struct {
		name          string
		pingError     error
		expectedError string
	}{
		{"Successful connection", nil, ""},
		{"Client version too old", fmt.Errorf("release_too_high"), "client version is too old"},
		{"Session evicted", fmt.Errorf("session evicted"), "session was evicted"},
		{"Generic error", fmt.Errorf("connection refused"), "connection failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Ping").Return(tt.pingError).Once()

			err := tb.ValidateConnectivity()

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if err == nil {
					t.Errorf("Expected an error containing '%s', but got nil", tt.expectedError)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			}

			// Check if the mock expectations were met
			expectationsMet := mockClient.AssertExpectations(t)
			if !expectationsMet {
				t.Errorf("Mock expectations were not met")
			}

			// Print the calls made to the mock
			calls := mockClient.Calls
			t.Logf("Number of calls made to mock: %d", len(calls))
			for i, call := range calls {
				t.Logf("Call %d: %s", i, call.Method)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	mockClient := new(MockClient)
	tb := &TigerBeagle{client: mockClient}

	// Test successful account creation
	mockClient.On("CreateAccounts", mock.Anything).Return(nil).Once()
	err := tb.CreateAccount(1)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)

	// Test failed account creation
	mockClient.On("CreateAccounts", mock.Anything).Return(fmt.Errorf("creation failed")).Once()
	err = tb.CreateAccount(2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creation failed")
	mockClient.AssertExpectations(t)
}

func TestTransfer(t *testing.T) {
	mockClient := new(MockClient)
	tb := &TigerBeagle{client: mockClient}

	// Test successful transfer
	mockClient.On("CreateTransfers", mock.Anything).Return(nil).Once()
	err := tb.Transfer(1, 2, 100)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)

	// Test failed transfer
	mockClient.On("CreateTransfers", mock.Anything).Return(fmt.Errorf("transfer failed")).Once()
	err = tb.Transfer(3, 4, 200)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transfer failed")
	mockClient.AssertExpectations(t)
}
