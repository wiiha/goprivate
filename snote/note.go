package snote

import (
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	// createdAt time.Time // `json:"createdAt"`
	// deletedAt time.Time // `json:"deletedAt"`
}

func NewNote(content string) *Note {
	id := nanoid.Must(15)
	return &Note{
		ID:      id,
		Content: content,
	}
}

/*
Equal is a utility function to
help when comparing two notes.
*/
func (n *Note) Equal(cn *Note) bool {
	return n.ID == cn.ID && n.Content == cn.Content
}

/*
NoteService works as a wrapper around
package specific functionality.

External calls are expected to go
through this service
*/
type NoteService struct {
	notesRepo *SQLiteRepository
}

/*
NewNoteService makes one
*/
func NewNoteService(repo *SQLiteRepository) *NoteService {
	return &NoteService{
		notesRepo: repo,
	}
}

func NewNoteServiceWithProductionRepo() *NoteService {
	repo := NewSqliteDefaultRepo()
	return NewNoteService(repo)
}

func (svc *NoteService) NewNote(content string) (string, error) {

	newNote := NewNote(content)

	err := svc.notesRepo.CreateNote(newNote)
	if err != nil {
		return "", fmt.Errorf("creating note: %v", err)
	}

	return newNote.ID, nil
}

/*
NotePingResponse will be the
response when calling PingNote
*/
type NotePingResponse struct {
	Exists    bool
	Deleted   bool
	DeletedAt time.Time
}

func (svc *NoteService) PingNote(noteID string) (*NotePingResponse, error) {

	deletedAt, err := svc.notesRepo.PingNote(noteID)
	if err != nil {
		return nil, fmt.Errorf("pinging note: %v", err)
	}

	if deletedAt == nil {
		return &NotePingResponse{
			Exists: false,
		}, nil
	}

	return &NotePingResponse{
		Exists:    true,
		Deleted:   deletedAt.Valid,
		DeletedAt: deletedAt.Time,
	}, nil
}

func (svc *NoteService) ConsumeNote(noteID string) (*Note, error) {

	note, err := svc.notesRepo.ConsumeNote(noteID)
	if err != nil {
		/*
			This can happen if a user
			requests a note that does
			not exists.
		*/
		if err == ErrNotExists {
			return nil, err
		}
		// Unexpected error here
		return nil, fmt.Errorf("consuming note: %v", err)
	}

	return note, nil
}
