package mapper

import (
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
)

func ToTradingJournalResponse(journal *entity.TradingJournal) *dto.TradingJournalResponse {
	return &dto.TradingJournalResponse{
		ID:          journal.ID,
		UserID:      journal.UserID,
		Name:        journal.Name,
		Description: journal.Description,
		CreatedAt:   journal.CreatedAt,
		UpdatedAt:   journal.UpdatedAt,
	}
}

func ToTradingJournalResponses(journals []*entity.TradingJournal) []*dto.TradingJournalResponse {
	responses := make([]*dto.TradingJournalResponse, len(journals))
	for i, journal := range journals {
		responses[i] = ToTradingJournalResponse(journal)
	}
	return responses
}

func ToTradingJournalWithEntriesResponse(journal *entity.TradingJournal) *dto.TradingJournalWithEntriesResponse {
	entries := make([]dto.TradingJournalEntryResponse, 0)
	if journal.Entries != nil {
		for _, entry := range journal.Entries {
			entries = append(entries, *ToTradingJournalEntryResponse(entry))
		}
	}

	return &dto.TradingJournalWithEntriesResponse{
		ID:          journal.ID,
		UserID:      journal.UserID,
		Name:        journal.Name,
		Description: journal.Description,
		Entries:     entries,
		CreatedAt:   journal.CreatedAt,
		UpdatedAt:   journal.UpdatedAt,
	}
}
