package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexkaplun/books-test/storage/models"
	"github.com/google/uuid"
)

func (s *storeImpl) CreateBook(ctx context.Context, book *models.Book) (*uuid.UUID, error) {
	var id uuid.UUID
	if err := s.db.QueryRowContext(ctx, createBook,
		book.Title, book.Author, book.Publisher, book.PublishDate, book.Rating, book.Status,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

// TODO: we may want to distinguish between errors if the id not found
func (s *storeImpl) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	if _, err := s.db.ExecContext(ctx, deleteBook, bookID); err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) UpdateBook(ctx context.Context, book *models.Book) error {
	res, err := s.db.ExecContext(ctx, updateBook,
		book.ID,
		book.Title,
		book.Author,
		book.Publisher,
		book.PublishDate,
		book.Rating,
		book.Status,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return ErrBookNotFound
	}

	return nil
}

func (s *storeImpl) GetBook(ctx context.Context, bookID uuid.UUID) (*models.Book, error) {
	var book models.Book
	if err := s.db.QueryRowContext(ctx, getBook, bookID).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Publisher,
		&book.PublishDate,
		&book.Rating,
		&book.Status,
		&book.CreatedAt,
		&book.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (s *storeImpl) ListBooks(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book
	rows, err := s.db.QueryContext(ctx, listBooks)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Publisher,
			&book.PublishDate,
			&book.Rating,
			&book.Status,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}
