package storage

import "context"

func (s *storeImpl) initialized(ctx context.Context) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(ctx, booksTableExists).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (s *storeImpl) createBooksTable(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, initSql); err != nil {
		return err
	}
	return nil
}
