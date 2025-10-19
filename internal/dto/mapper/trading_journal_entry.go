package mapper

import (
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
)

func ToTradingJournalEntryResponse(entry *entity.TradingJournalEntry) *dto.TradingJournalEntryResponse {
	return &dto.TradingJournalEntryResponse{
		ID:          entry.ID,
		JournalID:   entry.JournalID,
		Day:         entry.Day,
		Asset:       entry.Asset,
		LTF:         entry.LTF,
		HTF:         entry.HTF,
		EntryCharts: entry.EntryCharts,
		Session:     entry.Session,
		TradeType:   entry.TradeType,
		Setup:       entry.Setup,
		Direction:   entry.Direction,
		EntryType:   entry.EntryType,
		Realized:    entry.Realized,
		MaxRR:       entry.MaxRR,
		Result:      entry.Result,
		Notes:       entry.Notes,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
	}
}

func ToTradingJournalEntryResponses(entries []*entity.TradingJournalEntry) []*dto.TradingJournalEntryResponse {
	responses := make([]*dto.TradingJournalEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = ToTradingJournalEntryResponse(entry)
	}
	return responses
}

func ToStatisticsResponse(stats map[string]any) *dto.TradingJournalStatisticsResponse {
	response := &dto.TradingJournalStatisticsResponse{}

	if v, ok := stats["total_trades"].(int); ok {
		response.TotalTrades = v
	}
	if v, ok := stats["wins"].(int); ok {
		response.Wins = v
	}
	if v, ok := stats["losses"].(int); ok {
		response.Losses = v
	}
	if v, ok := stats["break_even"].(int); ok {
		response.BreakEven = v
	}
	if v, ok := stats["win_rate"].(float64); ok {
		response.WinRate = v
	}
	if v, ok := stats["total_realized"].(float64); ok {
		response.TotalRealized = v
	}
	if v, ok := stats["avg_risk_reward"].(float64); ok {
		response.AvgRiskReward = v
	}

	return response
}
