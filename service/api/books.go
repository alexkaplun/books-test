package api

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type UpsertBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publishDate"`
	Rating      int    `json:"rating"`
	Status      string `json:"status"`
}

type CreateBookResponse struct {
	ID *uuid.UUID `json:"id"`
}

func (m UpsertBookRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Author, validation.Required),
		validation.Field(&m.Rating, validation.Required, validation.Min(1), validation.Max(3)),
		validation.Field(&m.Status, validation.Required, validation.In("CheckedIn", "CheckedOut")),
		validation.Field(&m.PublishDate, validation.Date("2006-01-02")),
	)
}

type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate string    `json:"publishDate"`
	Rating      int       `json:"rating"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}
