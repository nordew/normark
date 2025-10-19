package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTradingJournalRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

type UpdateTradingJournalRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

type TradingJournalResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TradingJournalWithEntriesResponse struct {
	ID          uuid.UUID                    `json:"id"`
	UserID      uuid.UUID                    `json:"user_id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Entries     []TradingJournalEntryResponse `json:"entries"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
}

type TradingJournalListResponse struct {
	Journals []*TradingJournalResponse `json:"journals"`
	Total    int                       `json:"total"`
	Limit    int                       `json:"limit"`
	Offset   int                       `json:"offset"`
}
