package models

import (
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Transfer struct {
	ID              types.Uint128
	DebitAccountID  types.Uint128
	CreditAccountID types.Uint128
	Amount          types.Uint128
	PendingID       types.Uint128
	UserData128     types.Uint128
	UserData64      uint64
	UserData32      uint32
	Timeout         uint32
	Ledger          uint32
	Code            uint16
	Flags           uint16
	Timestamp       uint64
}

func (t *Transfer) ToTigerBeetleTransfer() types.Transfer {
	return types.Transfer{
		ID:              t.ID,
		DebitAccountID:  t.DebitAccountID,
		CreditAccountID: t.CreditAccountID,
		Amount:          t.Amount,
		PendingID:       t.PendingID,
		UserData128:     t.UserData128,
		UserData64:      t.UserData64,
		UserData32:      t.UserData32,
		Timeout:         t.Timeout,
		Ledger:          t.Ledger,
		Code:            t.Code,
		Flags:           t.Flags,
		Timestamp:       t.Timestamp,
	}
}

func FromTigerBeetleTransfer(tbt types.Transfer) *Transfer {
	return &Transfer{
		ID:              tbt.ID,
		DebitAccountID:  tbt.DebitAccountID,
		CreditAccountID: tbt.CreditAccountID,
		Amount:          tbt.Amount,
		PendingID:       tbt.PendingID,
		UserData128:     tbt.UserData128,
		UserData64:      tbt.UserData64,
		UserData32:      tbt.UserData32,
		Timeout:         tbt.Timeout,
		Ledger:          tbt.Ledger,
		Code:            tbt.Code,
		Flags:           tbt.Flags,
		Timestamp:       tbt.Timestamp,
	}
}
