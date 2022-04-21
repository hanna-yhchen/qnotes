package mysql

import (
	"database/sql"
	"errors"

	"github.com/hanna-yhchen/q-notes/internal/models"
)

type NoteModel struct {
	DB *sql.DB
}

// Insert inserts a new note into the database and returns the id for the new record.
func (m *NoteModel) Insert(userID int, title, content string) (int, error) {
	statement := `INSERT INTO notes (user_id, title, content, last_update) 
VALUES(?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(statement, userID, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get fetches a specific note by ID.
func (m *NoteModel) Get(id int) (*models.Note, error) {
	statement := `SELECT id, user_id, title, content, last_update FROM notes
WHERE id = ?`
	row := m.DB.QueryRow(statement, id)
	n := &models.Note{}

	if err := row.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.LastUpdate); err == nil {
		return n, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNoRecord
	} else {
		return nil, err
	}
}

// Latest returns the 10 most recently updated notes by a given user ID.
func (m *NoteModel) Latest(userID int) ([]*models.Note, error) {
	statement := `SELECT id, user_id, title, content, last_update FROM notes
WHERE user_id = ? ORDER BY last_update DESC LIMIT 10`

	rows, err := m.DB.Query(statement, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := []*models.Note{}

	for rows.Next() {
		n := &models.Note{}
		if err = rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.LastUpdate); err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// Update updates a specific note in the database.
func (m *NoteModel) Update(note *models.Note) error {
	statement := `UPDATE notes SET title = ?, content = ?, last_update = UTC_TIMESTAMP()
WHERE id = ?`

	_, err := m.DB.Exec(statement, note.Title, note.Content, note.ID)

	return err
}

// Delete deletes a specific note in the database.
func (m *NoteModel) Delete(id int) error {
	statement := `DELETE FROM notes WHERE id = ?`

	_, err := m.DB.Exec(statement, id)

	return err
}

// Search returns the 10 most recently updated notes where title or content contains the given string.
func (m *NoteModel) Search(text string) ([]*models.Note, error) {
	statement := `SELECT id, user_id, title, content, last_update FROM notes
WHERE MATCH (title,content) AGAINST(?) ORDER BY last_update DESC LIMIT 10`

	rows, err := m.DB.Query(statement, text)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return getNotesFromRows(rows)
}

// getNotesFromRows scans each row into a note object and returns the collection.
func getNotesFromRows(rows *sql.Rows) ([]*models.Note, error) {
	notes := []*models.Note{}

	for rows.Next() {
		n := &models.Note{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.LastUpdate); err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
