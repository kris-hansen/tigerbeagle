package tigerbeetle

import (
	"fmt"

	"github.com/kris-hansen/tigerbeagle/pkg/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Client interface {
	CreateAccounts(accounts []models.Account) error
	LookupAccount(id uint64) (*models.Account, error)
	CreateTransfers(transfers []models.Transfer) error
	Ping() error
	Close()
}

type tigerbeetleClient struct {
	client tb.Client
}

func NewClient(address string) (Client, error) {
	client, err := tb.NewClient(tbTypes.ToUint128(0), []string{address}, 256)
	if err != nil {
		return nil, fmt.Errorf("error creating TigerBeetle client: %w", err)
	}
	return &tigerbeetleClient{client: client}, nil
}

func (c *tigerbeetleClient) CreateAccounts(accounts []models.Account) error {
	tbAccounts := make([]tbTypes.Account, len(accounts))
	for i, account := range accounts {
		tbAccounts[i] = account.ToTigerBeetleAccount()
	}

	results, err := c.client.CreateAccounts(tbAccounts)
	if err != nil {
		return fmt.Errorf("error creating accounts: %w", err)
	}

	for _, result := range results {
		if result.Result != tbTypes.AccountOK {
			return fmt.Errorf("error creating account: %s", result.Result)
		}
	}

	return nil
}

func (c *tigerbeetleClient) LookupAccount(id uint64) (*models.Account, error) {
	accounts, err := c.client.LookupAccounts([]tbTypes.Uint128{tbTypes.ToUint128(id)})
	if err != nil {
		return nil, fmt.Errorf("error looking up account: %w", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("account not found")
	}

	return models.FromTigerBeetleAccount(accounts[0]), nil
}

func (c *tigerbeetleClient) CreateTransfers(transfers []models.Transfer) error {
	tbTransfers := make([]tbTypes.Transfer, len(transfers))
	for i, transfer := range transfers {
		tbTransfers[i] = transfer.ToTigerBeetleTransfer()
	}

	results, err := c.client.CreateTransfers(tbTransfers)
	if err != nil {
		return fmt.Errorf("error creating transfers: %w", err)
	}

	for _, result := range results {
		if result.Result != tbTypes.TransferOK {
			return fmt.Errorf("error creating transfer: %s", result.Result)
		}
	}

	return nil
}

func (c *tigerbeetleClient) Ping() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in Ping: %v", r)
		}
	}()

	_, err = c.client.LookupAccounts([]tbTypes.Uint128{tbTypes.ToUint128(0)})
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

func (c *tigerbeetleClient) Close() {
	c.client.Close()
}
