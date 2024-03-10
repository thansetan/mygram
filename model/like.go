package model

import "time"

type Like struct {
	ID        uint64
	UserID    uint64
	PhotoID   uint64
	CreatedAt time.Time

	User  User
	Photo Photo
}
