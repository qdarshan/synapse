package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteManager struct {
	DB *sql.DB
	saveNoteStmt *sql.Stmt
	deleteNoteStmt *sql.Stmt
}

type Note struct {
	Id              int
	Content         string
	EmbeddingVector []byte
	CreatedAt       time.Time
	Distance        float64
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func Initialize(filepath string) (*SQLiteManager, error) {

	db, err := sql.Open("sqlite_extended", filepath)

	if err != nil {
		logger.Error("Database: Failed to open connection file", "filepath", filepath, "error", err)
		return nil, err
	}

	logger.Debug("Database: Attempting to Ping connection...")
	if err := db.Ping(); err != nil {
		logger.Error("Database: Failed to Ping connection", "error", err)
		db.Close()
		return nil, err
	}
	logger.Debug("Database: Connection established successfully", "filepath", filepath)
	db.SetMaxOpenConns(1)

	manager := &SQLiteManager{DB: db}
	if err := manager.prepareStatements(); err != nil {
		db.Close()
		return nil, err
	}

	return manager, nil
}

func (manager *SQLiteManager) SetupSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    content TEXT NOT NULL,
	    embedding_vector BLOB NOT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	logger.Debug("Database: Setting up table schema...")
	_, err := manager.DB.Exec(schema)

	if err != nil {
		logger.Error("Database: Failed to execute schema creation", "error", err)
		return err
	}

	logger.Debug("Database: Schema created/verified successfully.")

	return nil
}

func (manager *SQLiteManager) prepareStatements() error {
	saveNoteQuery := `INSERT INTO notes (content, embedding_vector) VALUES(?, ?);`
	stmt, err := manager.DB.Prepare(saveNoteQuery)
	if err != nil {
		logger.Error("Database: Failed to prepare save note statement", "error", err)
		return fmt.Errorf("failed to prepare save note statement: %w", err)
	}
	manager.saveNoteStmt = stmt

	deleteNoteQuery := `DELETE FROM notes WHERE id = ?;`
	stmt, err = manager.DB.Prepare(deleteNoteQuery)
	if err != nil {
		logger.Error("Database: Failed to prepare delete note statement", "error", err)
		manager.saveNoteStmt.Close()
		return fmt.Errorf("failed to prepare delete note statement: %w", err)
	}
	manager.deleteNoteStmt = stmt

	logger.Debug("Database: Prepared statements cached successfully")
	return nil
}

func (manager *SQLiteManager) SaveNote(note Note) error {
	_, err := manager.saveNoteStmt.Exec(note.Content, note.EmbeddingVector)
	if err != nil {
		logger.Error("Database: Failed to EXECUTE statement for note insertion", "error", err)
		return fmt.Errorf("failed to execute statement for note insertion: %w", err)
	}

	logger.Debug("Database: Successfully saved a new note", "content_length", len(note.Content))

	return nil
}

func (manager *SQLiteManager) DeleteNote(id int) error {
	_, err := manager.deleteNoteStmt.Exec(id)
	if err != nil {
		logger.Error("Database: Failed to EXECUTE statement for note deletion", "error", err)
		return fmt.Errorf("failed to execute statement for note deletion: %w", err)
	}

	logger.Debug("Database: Successfully deleted note", "id", id)

	return nil
}

func (manager *SQLiteManager) GetNoteById(id int) (*Note, error) {
	getNoteByIdQuery := `SELECT id, content, embedding_vector, created_at FROM notes WHERE id = ?`

	rows, err := manager.DB.Query(getNoteByIdQuery, id)
	if err != nil {
		logger.Error("Database: Failed to execute SELECT query for note by id", "error", err)
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error during row iteration: %w", err)
		}
		return nil, nil
	}

	var note Note

	err = rows.Scan(
		&note.Id,
		&note.Content,
		&note.EmbeddingVector,
		&note.CreatedAt,
	)

	if err != nil {
		logger.Error("Database: Failed to scan row data into Note struct", "error", err)
		return nil, err
	}

	if rows.Next() {
		logger.Warn("Database: Query returned more than one note for a unique Id.", "id", note.Id)
	}

	return &note, nil
}

func (manager *SQLiteManager) GetAllNotes() ([]Note, error) {
	getAllNotesQuery := `SELECT id, content, embedding_vector, created_at FROM notes`

	rows, err := manager.DB.Query(getAllNotesQuery)
	if err != nil {
		logger.Error("Database: Failed to execute SELECT query for all notes", "error", err)
		return nil, err
	}

	notes := make([]Note, 0)
	defer rows.Close()

	for rows.Next() {
		var note Note
		err := rows.Scan(
			&note.Id,
			&note.Content,
			&note.EmbeddingVector,
			&note.CreatedAt,
		)

		if err != nil {
			logger.Error("Database: Failed to scan row data into Note struct", "error", err)
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Database: Error occurred during row iteration", "error", err)
		return nil, err
	}

	logger.Debug("Database: Successfully retrieved all notes", "count", len(notes))
	return notes, nil
}

func (manager *SQLiteManager) SearchNotes(queryVector []byte) ([]Note, error) {
	searchNotesQuery := `
	SELECT
	    id, content, embedding_vector, created_at,
	    vector_distance(embedding_vector, ?) AS distance
	FROM
	    notes
	ORDER BY
	    distance ASC
	LIMIT
	    10
	`

	rows, err := manager.DB.Query(searchNotesQuery, queryVector)
	if err != nil {
		logger.Error("Database: Failed to execute SELECT query for search notes", "error", err)
		return nil, err
	}

	notes := make([]Note, 0)
	defer rows.Close()

	for rows.Next() {
		var note Note
		var distance sql.NullFloat64
		err := rows.Scan(
			&note.Id,
			&note.Content,
			&note.EmbeddingVector,
			&note.CreatedAt,
			&distance,
		)

		if distance.Valid {
			note.Distance = distance.Float64
		} else {
			note.Distance = 999.0
		}

		if err != nil {
			logger.Error("Database: Failed to scan row data into Note struct", "error", err)
			return nil, err
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Database: Error occurred during row iteration", "error", err)
		return nil, err
	}

	logger.Debug("Database: Successfully found notes", "count", len(notes))
	return notes, nil
}
