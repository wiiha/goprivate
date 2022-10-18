package snote

import (
	"os"
	"testing"
)

/*
START Helpers
*/
func setupTestDB(t *testing.T) *SQLiteRepository {
	repo := NewSqliteTestRepo()

	err := repo.Migrate()
	if err != nil {
		t.Fatalf("migrating: %v", err)
	}
	return repo
}

func rmTestDB() {
	os.Remove("test.db")
}

/*
END Helpers
*/

func TestNoteServiceCreate(t *testing.T) {
	repo := setupTestDB(t)
	defer rmTestDB()

	svc := NewNoteService(repo)

	noteID, err := svc.NewNote("My test content")
	if err != nil {
		t.Fatalf("NewNote: %v", err)
	}

	if noteID == "" {
		t.Fatalf("expected to get id, got empty string")
	}

	// res, err := svc.PingNote(noteID)
	// if err != nil {
	// 	t.Fatalf("PingNote: %v", err)
	// }

	// if !res.Exists {
	// 	t.Fatalf("expected note with id %s to exist", noteID)
	// }

}

func TestNoteServicePing(t *testing.T) {
	repo := setupTestDB(t)
	defer rmTestDB()

	svc := NewNoteService(repo)

	res, err := svc.PingNote("cantfindme")
	if err != nil {
		t.Fatalf("PingNote: %v", err)
	}

	if res.Exists {
		t.Fatalf("expected note with id `cantfindme` to not exist, got %v", res)
	}

	noteID, err := svc.NewNote("My test content")
	if err != nil {
		t.Fatalf("NewNote: %v", err)
	}

	res, err = svc.PingNote(noteID)
	if err != nil {
		t.Fatalf("PingNote: %v", err)
	}

	if !res.Exists {
		t.Fatalf("expected note with id %s to exist, got %v", noteID, res)
	}

}

func TestNoteServiceConsume(t *testing.T) {
	repo := setupTestDB(t)
	defer rmTestDB()

	svc := NewNoteService(repo)

	content := "My test content"

	noteID, err := svc.NewNote(content)
	if err != nil {
		t.Fatalf("NewNote: %v", err)
	}

	ogNote := &Note{
		ID:      noteID,
		Content: content,
	}

	consumedNote, err := svc.ConsumeNote(noteID)
	if err != nil {
		t.Fatalf("ConsumeNote: %v", err)
	}

	if !consumedNote.Equal(ogNote) {
		t.Fatalf("consumed note note same as original. expected %v, got %v", ogNote, consumedNote)
	}
}
