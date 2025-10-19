package entity

import "github.com/cockroachdb/errors"

var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidJournalID   = errors.New("invalid journal ID")
	ErrInvalidJournalName = errors.New("invalid journal name")
	ErrInvalidAsset       = errors.New("invalid currency pair asset")
	ErrInvalidLTF         = errors.New("invalid lower timeframe (LTF)")
	ErrInvalidHTF         = errors.New("invalid higher timeframe (HTF)")
	ErrInvalidSession     = errors.New("invalid trading session")
	ErrInvalidTradeType   = errors.New("invalid trade type")
	ErrInvalidDirection   = errors.New("invalid trade direction")
	ErrInvalidEntryType   = errors.New("invalid entry type")
	ErrInvalidResult      = errors.New("invalid trade result")
)
