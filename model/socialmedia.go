package model

import "time"

type SocialMedia struct {
	ID, UserID           uint64
	Name, URL            string
	CreatedAt, UpdatedAt time.Time

	User User
}
