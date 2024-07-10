package models

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Account struct {
	ID             types.Uint128
	UserID         types.Uint128
	Ledger         uint32
	Code           uint16
	Flags          uint16
	DebitsPending  types.Uint128
	DebitsPosted   types.Uint128
	CreditsPending types.Uint128
	CreditsPosted  types.Uint128
}

// SetID sets the ID of the account using a uint64 value
func (a *Account) SetID(id uint64) {
	a.ID = types.ToUint128(id)
}

func (a Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID             uint64 `json:"id"`
		UserID         uint64 `json:"user_id"`
		Ledger         uint32 `json:"ledger"`
		Code           uint16 `json:"code"`
		Flags          uint16 `json:"flags"`
		DebitsPending  uint64 `json:"debits_pending"`
		DebitsPosted   uint64 `json:"debits_posted"`
		CreditsPending uint64 `json:"credits_pending"`
		CreditsPosted  uint64 `json:"credits_posted"`
	}{
		ID:             uint64FromUint128(a.ID),
		UserID:         uint64FromUint128(a.UserID),
		Ledger:         a.Ledger,
		Code:           a.Code,
		Flags:          a.Flags,
		DebitsPending:  uint64FromUint128(a.DebitsPending),
		DebitsPosted:   uint64FromUint128(a.DebitsPosted),
		CreditsPending: uint64FromUint128(a.CreditsPending),
		CreditsPosted:  uint64FromUint128(a.CreditsPosted),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *Account) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID             uint64 `json:"id"`
		UserID         uint64 `json:"user_id"`
		Ledger         uint32 `json:"ledger"`
		Code           uint16 `json:"code"`
		Flags          uint16 `json:"flags"`
		DebitsPending  uint64 `json:"debits_pending"`
		DebitsPosted   uint64 `json:"debits_posted"`
		CreditsPending uint64 `json:"credits_pending"`
		CreditsPosted  uint64 `json:"credits_posted"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.ID = types.ToUint128(aux.ID)
	a.UserID = types.ToUint128(aux.UserID)
	a.Ledger = aux.Ledger
	a.Code = aux.Code
	a.Flags = aux.Flags
	a.DebitsPending = types.ToUint128(aux.DebitsPending)
	a.DebitsPosted = types.ToUint128(aux.DebitsPosted)
	a.CreditsPending = types.ToUint128(aux.CreditsPending)
	a.CreditsPosted = types.ToUint128(aux.CreditsPosted)
	return nil
}

type AccountBalances struct {
	Account
	Balance *big.Int
}

func (a *Account) ToTigerBeetleAccount() types.Account {
	return types.Account{
		ID:             a.ID,
		UserData128:    a.UserID, // Assuming UserID maps to UserData128
		Ledger:         a.Ledger,
		Code:           a.Code,
		Flags:          a.Flags,
		DebitsPending:  a.DebitsPending,
		DebitsPosted:   a.DebitsPosted,
		CreditsPending: a.CreditsPending,
		CreditsPosted:  a.CreditsPosted,
	}
}

func (a *Account) FromTigerBeetleAccount(tba types.Account) {
	a.ID = tba.ID
	a.UserID = tba.UserData128 // Assuming UserData128 maps to UserID
	a.Ledger = tba.Ledger
	a.Code = tba.Code
	a.Flags = tba.Flags
	a.DebitsPending = tba.DebitsPending
	a.DebitsPosted = tba.DebitsPosted
	a.CreditsPending = tba.CreditsPending
	a.CreditsPosted = tba.CreditsPosted
}

func (ab *AccountBalances) FromTigerBeetleAccountBalance(tba types.AccountBalance) {
	ab.DebitsPending = tba.DebitsPending
	ab.DebitsPosted = tba.DebitsPosted
	ab.CreditsPending = tba.CreditsPending
	ab.CreditsPosted = tba.CreditsPosted
	// Note: We're not setting the Balance field here as it's not directly available in types.AccountBalance
}

func (ab *AccountBalances) ToTigerBeetleAccountBalances() (types.AccountBalance, error) {
	maxUint128 := types.ToUint128(^uint64(0))
	maxBigInt := maxUint128.BigInt()

	if ab.Balance.Cmp(&maxBigInt) > 0 {
		return types.AccountBalance{}, errors.New("balance exceeds maximum uint128 value")
	}

	return types.AccountBalance{
		DebitsPending:  ab.DebitsPending,
		DebitsPosted:   ab.DebitsPosted,
		CreditsPending: ab.CreditsPending,
		CreditsPosted:  ab.CreditsPosted,
		// Note: There's no direct field for Balance in types.AccountBalance
	}, nil
}

// Helper function to convert types.Uint128 to uint64
func uint64FromUint128(u types.Uint128) uint64 {
	bigInt := u.BigInt()
	return bigInt.Uint64()
}

func FromTigerBeetleAccount(tba types.Account) *Account {
	return &Account{
		ID:             tba.ID,
		UserID:         tba.UserData128, // Assuming UserData128 maps to UserID
		Ledger:         tba.Ledger,
		Code:           tba.Code,
		Flags:          tba.Flags,
		DebitsPending:  tba.DebitsPending,
		DebitsPosted:   tba.DebitsPosted,
		CreditsPending: tba.CreditsPending,
		CreditsPosted:  tba.CreditsPosted,
	}
}
