package snote

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wiiha/goprivate/database"
)

type SQLiteRepository struct {
	db *database.Database
	/*
		Since the repo uses a sqlite database we
		have to handle concurrency. The mutex will
		be used to ensure that only one process at
		a time accesses the DB file. This probably
		has performance implications.
	*/
	mutex sync.Mutex
}

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

func NewSqliteDefaultRepo() *SQLiteRepository {
	db := database.NewDatabase("sqlite3", "file:snote.db?_foreign_keys=true")
	return &SQLiteRepository{
		db:    db,
		mutex: sync.Mutex{},
	}
}

func NewSqliteTestRepo() *SQLiteRepository {
	db := database.NewDatabase("sqlite3", "file:test.db?_foreign_keys=true")
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS notes(
        id TEXT PRIMARY KEY,
        note_content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
    );
    `

	conn, err := r.db.Open()
	if err != nil {
		return fmt.Errorf("opening db conn: %v", err)
	}
	defer conn.Close()
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, err = conn.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateNote(n *Note) error {
	query := `
	INSERT INTO notes(id, note_content) values(?,?);
    `
	conn, err := r.db.Open()
	if err != nil {
		return fmt.Errorf("opening db conn: %v", err)
	}
	defer conn.Close()
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, err = conn.Exec(query, n.ID, n.Content)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}

/*
PingNote will check if a note
exists for the given noteID.

If the function returns without
error one should check the value of
the returned NullTime pointer.

+ If the pointer is nil then there
is no note with that id.

+ If NullTime.Valid is true then the note
exists and has been "deleted" i.e. read.

+ If NullTime.Valid is false then the note
exists and can still be read.
*/
func (r *SQLiteRepository) PingNote(noteID string) (*sql.NullTime, error) {
	conn, err := r.db.Open()
	if err != nil {
		return nil, fmt.Errorf("opening db conn: %v", err)
	}
	defer conn.Close()

	r.mutex.Lock()
	defer r.mutex.Unlock()
	row := conn.QueryRow("SELECT deleted_at FROM notes WHERE id = ?", noteID)

	deleted_at := sql.NullTime{}
	if err := row.Scan(&deleted_at); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying: %v", err)
	}
	return &deleted_at, nil
}

func (r *SQLiteRepository) ConsumeNote(noteID string) (*Note, error) {
	conn, err := r.db.Open()
	if err != nil {
		return nil, fmt.Errorf("opening db conn: %v", err)
	}
	defer conn.Close()

	r.mutex.Lock()
	defer r.mutex.Unlock()
	row := conn.QueryRow("SELECT id, note_content FROM notes WHERE id = ?", noteID)

	var note Note
	if err := row.Scan(&note.ID, &note.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, fmt.Errorf("querying: %v", err)
	}

	/*
		var `note` contains the note that
		we should return. Time to remove
		the content from the DB.
	*/

	query := `
	UPDATE notes 
	SET note_content = NULL, 
	deleted_at = CURRENT_TIMESTAMP
	WHERE id = ?
    `
	res, err := conn.Exec(query, noteID)
	if err != nil {
		return nil, fmt.Errorf("execing: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("RowsAffected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &note, nil
}
