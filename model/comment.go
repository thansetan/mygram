package model

import "time"

type Comment struct {
	ID, UserID, PhotoID  uint64
	Message              string
	CreatedAt, UpdatedAt time.Time

	User  User
	Photo Photo
}
