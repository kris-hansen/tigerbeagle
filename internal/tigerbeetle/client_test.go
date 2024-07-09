package tigerbeetle

import (
	"fmt"
	"testing"

	"github.com/kris-hansen/tigerbeagle/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Mock tb.Client
type MockTBClient struct {
	mock.Mock
}

func (m *MockTBClient) CreateAccounts(accounts []tbTypes.Account) ([]tbTypes.AccountEventResult, error) {
	args := m.Called(accounts)
	return args.Get(0).([]tbTypes.AccountEventResult), args.Error(1)
}

func (m *MockTBClient) LookupAccounts(ids []tbTypes.Uint128) ([]tbTypes.Account, error) {
	args := m.Called(ids)
	return args.Get(0).([]tbTypes.Account), args.Error(1)
}

func (m *MockTBClient) CreateTransfers(transfers []tbTypes.Transfer) ([]tbTypes.TransferEventResult, error) {
	args := m.Called(transfers)
	return args.Get(0).([]tbTypes.TransferEventResult), args.Error(1)
}

func (m *MockTBClient) LookupTransfers(ids []tbTypes.Uint128) ([]tbTypes.Transfer, error) {
	args := m.Called(ids)
	return args.Get(0).([]tbTypes.Transfer), args.Error(1)
}

func (m *MockTBClient) GetAccountTransfers(filter tbTypes.AccountFilter) ([]tbTypes.Transfer, error) {
	args := m.Called(filter)
	return args.Get(0).([]tbTypes.Transfer), args.Error(1)
}

func (m *MockTBClient) GetAccountBalances(filter tbTypes.AccountFilter) ([]tbTypes.AccountBalance, error) {
	args := m.Called(filter)
	return args.Get(0).([]tbTypes.AccountBalance), args.Error(1)
}

func (m *MockTBClient) Nop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTBClient) Close() {
	m.Called()
}
func TestCreateAccounts(t *testing.T) {
	mockTB := new(MockTBClient)
	client := &tigerbeetleClient{client: mockTB}

	// Test successful account creation
	mockTB.On("CreateAccounts", mock.Anything).Return([]tbTypes.AccountEventResult{{Result: tbTypes.AccountOK}}, nil).Once()
	err := client.CreateAccounts([]models.Account{{ID: tbTypes.ToUint128(1)}})
	assert.NoError(t, err)
	mockTB.AssertExpectations(t)

	// Test account creation failure
	mockTB.On("CreateAccounts", mock.Anything).Return([]tbTypes.AccountEventResult{{Result: tbTypes.AccountExists}}, nil).Once()
	err = client.CreateAccounts([]models.Account{{ID: tbTypes.ToUint128(2)}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating account: AccountExists")
	mockTB.AssertExpectations(t)
}

func TestLookupAccount(t *testing.T) {
	mockTB := new(MockTBClient)
	client := &tigerbeetleClient{client: mockTB}

	// Test successful account lookup
	mockTB.On("LookupAccounts", mock.Anything).Return([]tbTypes.Account{{ID: tbTypes.ToUint128(1)}}, nil).Once()
	account, err := client.LookupAccount(1)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	mockTB.AssertExpectations(t)

	// Test account not found
	mockTB.On("LookupAccounts", mock.Anything).Return([]tbTypes.Account{}, nil).Once()
	account, err = client.LookupAccount(2)
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "account not found")
	mockTB.AssertExpectations(t)

	// Test lookup error
	mockTB.On("LookupAccounts", mock.Anything).Return([]tbTypes.Account{}, fmt.Errorf("lookup error")).Once()
	account, err = client.LookupAccount(3)
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "error looking up account: lookup error")
	mockTB.AssertExpectations(t)
}
