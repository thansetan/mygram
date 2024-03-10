package model

import (
	"database/sql"
	"time"
)

type Photo struct {
	ID, UserID           uint64
	Title, URL           string
	Caption              sql.NullString
	CreatedAt, UpdatedAt time.Time

	User     User
	Comments []Comment
}
