package storage

const (
	initSql = `
CREATE TABLE books (
	id				UUID			DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
	title			VARCHAR(255)	NOT NULL,
	author			VARCHAR(255)	NOT NULL,
	publisher		VARCHAR(255)	NULL,
	publish_date	DATE			NULL,
	rating			INTEGER			NULL,
	status			VARCHAR(64)		NOT NULL,

	created_at		TIMESTAMP		NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at		TIMESTAMP		NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

	booksTableExists = `
SELECT EXISTS (
   SELECT FROM information_schema.tables 
   WHERE  table_schema = 'public'
   AND    table_name   = 'books'
   );
`

	createBook = `
INSERT INTO books
	(title, author, publisher, publish_date, rating, status)
VALUES 
	($1, $2, $3, $4, $5, $6)
RETURNING
	id
`

	deleteBook = `
DELETE FROM books
WHERE id = $1
`

	updateBook = `
UPDATE books
SET title = $2, author = $3, publisher = $4, publish_date = $5, rating = $6, status = $7,
	updated_at = CURRENT_TIMESTAMP
WHERE 
	id = $1
`

	getBook = `
SELECT 
	id, title, author, publisher, publish_date, rating, status, created_at, updated_at
FROM books
WHERE id = $1
`

	listBooks = `
SELECT 
	id, title, author, publisher, publish_date, rating, status, created_at, updated_at
FROM books
ORDER BY created_at DESC
`
)
