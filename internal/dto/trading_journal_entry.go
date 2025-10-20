package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user/normark/internal/types"
)

type CreateTradingJournalEntryRequest struct {
	Day         time.Time              `json:"day" validate:"required"`
	Asset       types.CurrencyPair     `json:"asset" validate:"required"`
	LTF         string                 `json:"ltf" validate:"required,url"`
	HTF         string                 `json:"htf" validate:"required,url"`
	EntryCharts []string               `json:"entry_charts" validate:"omitempty,dive,url"`
	Session     types.TradingSession   `json:"session" validate:"required"`
	TradeType   types.TradeType        `json:"trade_type" validate:"required"`
	Setup       *string                `json:"setup" validate:"omitempty,max=500"`
	Direction   types.TradeDirection   `json:"direction" validate:"required"`
	EntryType   types.EntryType        `json:"entry_type" validate:"required"`
	Realized    float64                `json:"realized" validate:"required"`
	MaxRR       float64                `json:"max_rr" validate:"required,gt=0"`
	Result      types.TradeResult      `json:"result" validate:"required"`
	Notes       string                 `json:"notes" validate:"omitempty,max=5000"`
}

type UpdateTradingJournalEntryRequest struct {
	Day         time.Time              `json:"day" validate:"required"`
	Asset       types.CurrencyPair     `json:"asset" validate:"required"`
	LTF         string                 `json:"ltf" validate:"required,url"`
	HTF         string                 `json:"htf" validate:"required,url"`
	EntryCharts []string               `json:"entry_charts" validate:"omitempty,dive,url"`
	Session     types.TradingSession   `json:"session" validate:"required"`
	TradeType   types.TradeType        `json:"trade_type" validate:"required"`
	Setup       *string                `json:"setup" validate:"omitempty,max=500"`
	Direction   types.TradeDirection   `json:"direction" validate:"required"`
	EntryType   types.EntryType        `json:"entry_type" validate:"required"`
	Realized    float64                `json:"realized" validate:"required"`
	MaxRR       float64                `json:"max_rr" validate:"required,gt=0"`
	Result      types.TradeResult      `json:"result" validate:"required"`
	Notes       string                 `json:"notes" validate:"omitempty,max=5000"`
}

type TradingJournalEntryResponse struct {
	ID          uuid.UUID              `json:"id"`
	JournalID   uuid.UUID              `json:"journal_id"`
	Day         time.Time              `json:"day"`
	Asset       types.CurrencyPair     `json:"asset"`
	LTF         string                 `json:"ltf"`
	HTF         string                 `json:"htf"`
	EntryCharts []string               `json:"entry_charts"`
	Session     types.TradingSession   `json:"session"`
	TradeType   types.TradeType        `json:"trade_type"`
	Setup       *string                `json:"setup,omitempty"`
	Direction   types.TradeDirection   `json:"direction"`
	EntryType   types.EntryType        `json:"entry_type"`
	Realized    float64                `json:"realized"`
	MaxRR       float64                `json:"max_rr"`
	Result      types.TradeResult      `json:"result"`
	Notes       string                 `json:"notes"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TradingJournalEntryListResponse struct {
	Entries []*TradingJournalEntryResponse `json:"entries"`
	Total   int                            `json:"total"`
	Limit   int                            `json:"limit"`
	Offset  int                            `json:"offset"`
}

type TradingJournalStatisticsResponse struct {
	TotalTrades     int     `json:"total_trades"`
	Wins            int     `json:"wins"`
	Losses          int     `json:"losses"`
	BreakEven       int     `json:"break_even"`
	WinRate         float64 `json:"win_rate"`
	TotalRealized   float64 `json:"total_realized"`
	AvgRiskReward   float64 `json:"avg_risk_reward"`
}

type FilterEntriesRequest struct {
	Asset     *types.CurrencyPair   `json:"asset" validate:"omitempty"`
	Session   *types.TradingSession `json:"session" validate:"omitempty"`
	Result    *types.TradeResult    `json:"result" validate:"omitempty"`
	StartDate *time.Time            `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time            `json:"end_date" validate:"omitempty"`
	Limit     int                   `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int                   `json:"offset" validate:"omitempty,min=0"`
}
