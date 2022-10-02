package snote

import (
	"os"
	"testing"
)

func TestMigrate(t *testing.T) {
	repo := NewSqliteTestRepo()

	err := repo.Migrate()
	if err != nil {
		t.Fatalf("migrating: %v", err)
	}
	defer os.Remove("test.db")
}

func TestCreateAndPing(t *testing.T) {
	repo := NewSqliteTestRepo()

	err := repo.Migrate()
	if err != nil {
		t.Fatalf("migrating: %v", err)
	}
	defer os.Remove("test.db")

	n1 := &Note{
		ID:      "test1",
		Content: "this is test content",
	}

	err = repo.CreateNote(n1)
	if err != nil {
		t.Fatalf("creating note: %v", err)
	}

	deleted_at, err := repo.PingNote(n1.ID)
	if err != nil {
		t.Fatalf("pinging note: %v", err)
	}
	if deleted_at == nil {
		t.Fatalf("Expected to find note with id %s, got nil", n1.ID)
	}

	if deleted_at.Valid {
		t.Fatalf("Expected note to still be read able, got %v", deleted_at)
	}

	deleted_at, err = repo.PingNote("cannotfindme")
	if err != nil {
		t.Fatalf("pinging note: %v", err)
	}
	if deleted_at != nil {
		t.Fatalf("Did not expect to find a note but got true on ping")
	}
}

func TestConsumeNote(t *testing.T) {
	repo := NewSqliteTestRepo()

	err := repo.Migrate()
	if err != nil {
		t.Fatalf("migrating: %v", err)
	}
	defer os.Remove("test.db")

	n1 := &Note{
		ID:      "test1",
		Content: "this is test content",
	}

	err = repo.CreateNote(n1)
	if err != nil {
		t.Fatalf("creating note: %v", err)
	}

	retNote, err := repo.ConsumeNote(n1.ID)
	if err != nil {
		t.Fatalf("consuming note: %v", err)
	}

	if n1.Content != retNote.Content {
		t.Fatalf("expected content \"%s\", got \"%s\"", n1.Content, retNote.Content)
	}

	deleted_at, err := repo.PingNote(n1.ID)
	if err != nil {
		t.Fatalf("pinging note: %v", err)
	}
	if deleted_at == nil {
		t.Fatalf("Expected to find note with id %s, got nil", n1.ID)
	}

	if !deleted_at.Valid {
		t.Fatalf("expected note to be marked as deleted, got %v", deleted_at)
	}
}
