package models

import (
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

// Add this function
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
