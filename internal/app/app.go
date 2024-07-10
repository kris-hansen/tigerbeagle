package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kris-hansen/tigerbeagle/internal/tigerbeetle"
	"github.com/kris-hansen/tigerbeagle/pkg/models"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBeagleInterface interface {
	ValidateConnectivity() error
	CreateAccount(id uint64, ledger uint32, code uint16, flags uint16) error
	GetAccount(id uint64) (*models.Account, error)
	Transfer(debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error
	BulkTransfer(iterations int, debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error
	GenerateAccounts(number int, ledger uint32, code uint16, flags uint16) error
	GenerateTransfers(number int, ledger uint32, code uint16, flags uint16) error
	MigrateAccounts(filename string) error
	MigrateTransfers(filename string) error
}

var _ TigerBeagleInterface = (*TigerBeagle)(nil)

type TigerBeagle struct {
	client tigerbeetle.Client
}

func NewTigerBeagle() *TigerBeagle {
	return &TigerBeagle{}
}

func (t *TigerBeagle) InitClient(address string) error {
	client, err := tigerbeetle.NewClient(address)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	t.client = client
	return nil
}

func (t *TigerBeagle) CloseClient() {
	if t.client != nil {
		t.client.Close()
	}
}

func (t *TigerBeagle) CreateAccount(id uint64, ledger uint32, code uint16, flags uint16) error {
	account := models.Account{
		ID:             tbTypes.ToUint128(id),
		DebitsPending:  tbTypes.ToUint128(0),
		DebitsPosted:   tbTypes.ToUint128(0),
		CreditsPending: tbTypes.ToUint128(0),
		CreditsPosted:  tbTypes.ToUint128(0),
		UserID:         tbTypes.ToUint128(0),
		Ledger:         ledger,
		Code:           code,
		Flags:          flags,
	}

	err := t.client.CreateAccounts([]models.Account{account})
	if err != nil {
		return fmt.Errorf("error creating account: %w", err)
	}

	fmt.Printf("Account created with ID: %d, Ledger: %d, Code: %d, Flags: %d\n", id, ledger, code, flags)
	return nil
}

func (t *TigerBeagle) GetAccount(id uint64) (*models.Account, error) {
	account, err := t.client.LookupAccount(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching account: %w", err)
	}
	return account, nil
}

func (t *TigerBeagle) Transfer(debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error {
	transfer := models.Transfer{
		ID:              tbTypes.ToUint128(uint64(time.Now().UnixNano())),
		DebitAccountID:  tbTypes.ToUint128(debitAccountID),
		CreditAccountID: tbTypes.ToUint128(creditAccountID),
		Amount:          tbTypes.ToUint128(amount),
		Ledger:          ledger,
		Code:            code,
		Flags:           flags,
	}

	err := t.client.CreateTransfers([]models.Transfer{transfer})
	if err != nil {
		return fmt.Errorf("error creating transfer: %w", err)
	}

	fmt.Printf("Transfer completed: %d from account %d to account %d (Ledger: %d, Code: %d, Flags: %d)\n",
		amount, debitAccountID, creditAccountID, ledger, code, flags)
	return nil
}

func (t *TigerBeagle) BulkTransfer(iterations int, debitAccountID, creditAccountID, amount uint64, ledger uint32, code uint16, flags uint16) error {
	for i := 0; i < iterations; i++ {
		transfer := models.Transfer{
			ID:              tbTypes.ToUint128(uint64(time.Now().UnixNano()) + uint64(i)),
			DebitAccountID:  tbTypes.ToUint128(debitAccountID),
			CreditAccountID: tbTypes.ToUint128(creditAccountID),
			Amount:          tbTypes.ToUint128(amount),
			Ledger:          ledger,
			Code:            code,
			Flags:           flags,
		}

		err := t.client.CreateTransfers([]models.Transfer{transfer})
		if err != nil {
			return fmt.Errorf("error creating transfer in iteration %d: %w", i, err)
		}

		fmt.Printf("Transfer %d completed: %d from account %d to account %d (Ledger: %d, Code: %d, Flags: %d)\n",
			i+1, amount, debitAccountID, creditAccountID, ledger, code, flags)
	}

	return nil
}

func (t *TigerBeagle) MigrateAccounts(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var accounts []models.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	err = t.client.CreateAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error creating accounts: %w", err)
	}

	fmt.Printf("Successfully migrated %d accounts\n", len(accounts))
	return nil
}

func (t *TigerBeagle) MigrateTransfers(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var transfers []models.Transfer
	if err := json.Unmarshal(data, &transfers); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	err = t.client.CreateTransfers(transfers)
	if err != nil {
		return fmt.Errorf("error creating transfers: %w", err)
	}

	fmt.Printf("Successfully migrated %d transfers\n", len(transfers))
	return nil
}

func (t *TigerBeagle) ValidateConnectivity() error {
	err := t.client.Ping()

	if err != nil {
		if strings.Contains(err.Error(), "release_too_high") {
			return fmt.Errorf("client version is too old: %w", err)
		}
		if strings.Contains(err.Error(), "session evicted") {
			return fmt.Errorf("session was evicted: %w", err)
		}
		return fmt.Errorf("connection failed: %w", err)
	}

	return nil
}

func (t *TigerBeagle) GenerateAccounts(number int, ledger uint32, code uint16, flags uint16) error {
	accounts := make([]models.Account, number)
	for i := 0; i < number; i++ {
		account := models.Account{
			UserID:         types.ToUint128(0),
			Ledger:         ledger,
			Code:           code,
			Flags:          flags,
			DebitsPending:  types.ToUint128(0),
			DebitsPosted:   types.ToUint128(0),
			CreditsPending: types.ToUint128(0),
			CreditsPosted:  types.ToUint128(0),
		}
		account.SetID(uint64(1000 + i))
		accounts[i] = account
	}

	return writeJSONToFile(accounts, "generated_accounts.json")
}

func (t *TigerBeagle) GenerateTransfers(number int, ledger uint32, code uint16, flags uint16) error {
	transfers := make([]models.Transfer, number)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < number; i++ {
		transfers[i] = models.Transfer{
			ID:              tbTypes.ToUint128(uint64(i + 1)),
			DebitAccountID:  tbTypes.ToUint128(uint64(rand.Intn(number) + 1000)),
			CreditAccountID: tbTypes.ToUint128(uint64(rand.Intn(number) + 1000)),
			Amount:          tbTypes.ToUint128(uint64(rand.Intn(10000) + 1)),
			Ledger:          ledger,
			Code:            code,
			Flags:           flags,
		}
	}

	return writeJSONToFile(transfers, "generated_transfers.json")
}

func writeJSONToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	fmt.Printf("Generated %s successfully.\n", filename)
	return nil
}
