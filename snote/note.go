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

func newNote(content string) *Note {
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

/*
NewNote will create a note with the given content
and stor this note in the datastore.

The function returns the ID of the newly created note.
*/
func (svc *NoteService) NewNote(content string) (string, error) {

	newNote := newNote(content)

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
	Exists    bool      `json:"exists"`
	Deleted   bool      `json:"consumed"`
	DeletedAt time.Time `json:"consumedAt"`
}

/*
PingNote checks if the note with the
given ID exists and what "state" the
note is in.

States:

1. Does not exist

2. Exists and unread

3. Exists and read
*/
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

/*
ConsumeNote will try retrieving the note
with the given ID. If it is successful in
fetching the note it will also be deleted
from the datastore. Meaning that when the
note is returned from this function it does
no longer exist in the database and cannot
be retrieved agian.
*/
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
