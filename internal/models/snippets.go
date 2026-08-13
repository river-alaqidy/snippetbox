package models

import (
	"database/sql"
	"errors"
	"time"
)

// snippet type to hold data for invidual snippet
// corresponds to fields in talble
type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// wrap db connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// TODO: look into golan embed, can put sql statements in .sql files for cleanliness
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// last 3 params are placeholders that replace the '?'s in sql statement
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get id of newly created snippet record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// id returned is of type int64, needs to be converted
	return int(id), err
}

// Get a snippet based on id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	// copy values from each field to snippet struct
	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if query returns no rows
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Latest returns 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	// iterate through results
	for rows.Next() {
		var s Snippet
		// copy values from each field in the row to new snippet
		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// retrieve any error that occurred during iteration after completion
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
