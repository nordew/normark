package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/user/normark/internal/types"
)

type TradingJournalEntry struct {
	bun.BaseModel `bun:"table:trading_journal_entries,alias:tje"`

	ID          uuid.UUID            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	JournalID   uuid.UUID            `bun:"journal_id,notnull,type:uuid"`
	Day         time.Time            `bun:"day,notnull"`
	Asset       types.CurrencyPair   `bun:"asset,notnull"`
	LTF         types.TimeFrame      `bun:"ltf,notnull"`
	HTF         types.TimeFrame      `bun:"htf,notnull"`
	EntryCharts []string             `bun:"entry_charts,array,type:text[]"`
	Session     types.TradingSession `bun:"session,notnull"`
	TradeType   types.TradeType      `bun:"trade_type,notnull"`
	Setup       *string              `bun:"setup,nullzero"`
	Direction   types.TradeDirection `bun:"direction,notnull"`
	EntryType   types.EntryType      `bun:"entry_type,notnull"`
	Realized    float64              `bun:"realized,type:decimal(10,2),notnull"`
	MaxRR       float64              `bun:"max_rr,type:decimal(10,2),notnull"`
	Result      types.TradeResult    `bun:"result,notnull"`
	Notes       string               `bun:"notes,type:text"`
	CreatedAt   time.Time            `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time            `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt   time.Time            `bun:"deleted_at,soft_delete,nullzero"`

	Journal *TradingJournal `bun:"rel:belongs-to,join:journal_id=id"`
}

func NewTradingJournalEntry(
	journalID uuid.UUID,
	day time.Time,
	asset types.CurrencyPair,
	ltf, htf types.TimeFrame,
	entryCharts []string,
	session types.TradingSession,
	tradeType types.TradeType,
	setup *string,
	direction types.TradeDirection,
	entryType types.EntryType,
	realized, maxRR float64,
	result types.TradeResult,
	notes string,
) *TradingJournalEntry {
	return &TradingJournalEntry{
		JournalID:   journalID,
		Day:         day,
		Asset:       asset,
		LTF:         ltf,
		HTF:         htf,
		EntryCharts: entryCharts,
		Session:     session,
		TradeType:   tradeType,
		Setup:       setup,
		Direction:   direction,
		EntryType:   entryType,
		Realized:    realized,
		MaxRR:       maxRR,
		Result:      result,
		Notes:       notes,
	}
}

func (tje *TradingJournalEntry) Validate() error {
	if tje.JournalID == uuid.Nil {
		return ErrInvalidJournalID
	}

	if !tje.Asset.IsValid() {
		return ErrInvalidAsset
	}

	if !tje.LTF.IsValid() {
		return ErrInvalidLTF
	}

	if !tje.HTF.IsValid() {
		return ErrInvalidHTF
	}

	if !tje.Session.IsValid() {
		return ErrInvalidSession
	}

	if !tje.TradeType.IsValid() {
		return ErrInvalidTradeType
	}

	if !tje.Direction.IsValid() {
		return ErrInvalidDirection
	}

	if !tje.EntryType.IsValid() {
		return ErrInvalidEntryType
	}

	if !tje.Result.IsValid() {
		return ErrInvalidResult
	}

	return nil
}

func (tje *TradingJournalEntry) IsProfit() bool {
	return tje.Realized > 0
}

func (tje *TradingJournalEntry) IsLoss() bool {
	return tje.Realized < 0
}

func (tje *TradingJournalEntry) IsBreakEven() bool {
	return tje.Realized == 0
}
