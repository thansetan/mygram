package model

import "time"

type User struct {
	ID, Age              uint64
	Username, Email      string
	Password             []byte
	CreatedAt, UpdatedAt time.Time
}
