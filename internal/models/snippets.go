package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
			AND id = ?
	`

	// Returns a pointer to a sql.Row object which
	// holds the result from the database
	row := m.DB.QueryRow(stmt, id)

	var s Snippet
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct
	// The number of arguments must be exactly the same as the number of
	// columns returned by your statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// This helps to encapsulate the model completely,
			// so that our handlers aren’t concerned with the
			// underlying datastore or reliant on datastore-specific
			// errors
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
