package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Publisher   string
	PublishDate *time.Time
	Rating      int
	Status      BookStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type BookStatus string

const (
	BookStatusCheckedIn  BookStatus = "CheckedIn"
	BookStatusCheckedOut BookStatus = "CheckedOut"
)
