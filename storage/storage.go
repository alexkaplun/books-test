package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/alexkaplun/books-test/storage/models"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateBook(ctx context.Context, book *models.Book) (*uuid.UUID, error)
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
	UpdateBook(ctx context.Context, book *models.Book) error
	GetBook(ctx context.Context, bookID uuid.UUID) (*models.Book, error)
	ListBooks(ctx context.Context) ([]*models.Book, error)
}

type Params struct {
	ConnString string
}

type storeImpl struct {
	db *sql.DB
}

func NewPostgres(params Params) (Storage, error) {
	db, err := sql.Open("postgres", params.ConnString)
	if err != nil {
		return nil, err
	}

	store := &storeImpl{
		db: db,
	}

	if err = store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *storeImpl) init() error {
	// allow up to 5 seconds to init the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tableExists, err := s.initialized(ctx)
	if err != nil {
		return err
	}

	if !tableExists {
		if err = s.createBooksTable(ctx); err != nil {
			return err
		}
	}

	return nil
}
