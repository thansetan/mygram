package dto

import "time"

type LikeRequest struct {
	PhotoID uint64 `json:"-"`
}

type LikeCreateResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"user"`
}

type GetLikeByUserIDResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`

	Photo Photo `json:"photo"`
}
