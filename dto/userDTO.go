package dto

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	Page       int            `json:"page"`
	PrevPage   int            `json:"prev_page"`
	NextPage   int            `json:"next_page"`
	HasPrev    bool           `json:"has_prev"`
	HasNext    bool           `json:"has_next"`
	TotalPages int            `json:"total_pages"`
	Search     string         `json:"search"`
}
