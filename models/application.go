package models

import "time"

type Application struct {
	ID        uint
	UserID    uint
	Text      string
	FileURL   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateApplicationRequest struct {
	Text    string `json:"text" binding:"required"`
	Status  string `json:"status"`
	FileURL string `json:"file_url"  binding:"required"`
}

type UpdateApplicationRequest struct {
	Text    string `json:"text"`
	Status  string `json:"status"`
	FileURL string `json:"file_url"`
}
