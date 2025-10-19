package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TradingJournal struct {
	bun.BaseModel `bun:"table:trading_journals,alias:tj"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID      uuid.UUID `bun:"user_id,notnull,type:uuid"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description,type:text"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt   time.Time `bun:"deleted_at,soft_delete,nullzero"`

	User    *User                  `bun:"rel:belongs-to,join:user_id=id"`
	Entries []*TradingJournalEntry `bun:"rel:has-many,join:id=journal_id"`
}

func NewTradingJournal(userID uuid.UUID, name, description string) *TradingJournal {
	return &TradingJournal{
		UserID:      userID,
		Name:        name,
		Description: description,
	}
}

func (tj *TradingJournal) Validate() error {
	if tj.UserID == uuid.Nil {
		return ErrInvalidUserID
	}

	if tj.Name == "" {
		return ErrInvalidJournalName
	}

	return nil
}
