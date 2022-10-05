package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wiiha/goprivate/snote"
)

type Server struct {
	httpServer *http.Server
}

/*
NewServer makes a server. The function
handles creating a router and registering
routes to the router.

The returned Server struct has a ListenAndServe
method which should be invoked in order for
the server to start.
*/
func NewServer(addressAndPort string, noteSVC *snote.NoteService) *Server {
	router := gin.Default()

	/***
		Register routes to router
	***/

	/*
		The landing page must be specified as
		a file in order to not create a
		`/*filepath` wildcard route that would
		conflict with all other routes.
	*/
	router.StaticFileFS("/", "./", fsLandigPage())

	/*
		The assets folder contains multiple files
		and is therefore served as a FS rather than
		a specific set of files. This works since
		the registered route will be `/assets/*filepath`.
	*/
	router.StaticFS("/assets", fsAssetsDir())

	/*
		This route is used to read a message. Actual fetching
		of the message will be handled by the frontend.
		This route exists in order to emulate a single
		page application.
	*/
	router.GET("/read/:messageid", func(c *gin.Context) {
		messageid := c.Param("messageid")
		log.Printf("(debug) messageid: %v", messageid)

		c.FileFromFS("./", fsLandigPage())
	})

	/*
		API for backend
	*/

	/*
		Check if API is alive
	*/
	router.GET("api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	/*
		Create a new note. The received content is already
		expected to be encrypted.
	*/
	router.POST("api/v1/newnote", func(c *gin.Context) {
		var newNote NewNote
		if err := c.ShouldBindJSON(&newNote); err != nil {
			c.String(http.StatusBadRequest, "field noteContent is invalid")
			return
		}

		noteID, err := noteSVC.NewNote(newNote.NoteContent)
		if err != nil {
			log.Printf("(error) api/v1/newnote: %v", err)
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, gin.H{"noteID": noteID})
	})

	/*
		Used to check if a note exists without fetching
		the actual note.
	*/
	router.GET("api/v1/pingnote/:noteid", func(c *gin.Context) {
		noteID := c.Param("noteid")

		pingRes, err := noteSVC.PingNote(noteID)
		if err != nil {
			log.Printf("(error) api/v1/pingnote/%s: %v", noteID, err)
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, pingRes)
	})

	/*
		When a note is returned by this route it has been
		removed from the database and there is no way to
		recover the previous content of the note.
	*/
	router.DELETE("api/v1/consumenote/:noteid", func(c *gin.Context) {
		noteID := c.Param("noteid")

		consumedNote, err := noteSVC.ConsumeNote(noteID)
		if err != nil {
			if err == snote.ErrNotExists {
				c.String(http.StatusBadRequest, "note does not exist")
				return
			}
			log.Printf("(error) api/v1/consume/%s: %v", noteID, err)
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, consumedNote)
	})

	/***
		END register routes
	***/

	s := &http.Server{
		Addr:         addressAndPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		httpServer: s,
	}
}

func (s *Server) ListenAndServe() error {

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}
	return nil
}

/*
struct NewNote is used for binding to the
json post request representing a new note.
*/
type NewNote struct {
	NoteContent string `json:"noteContent" binding:"required"`
}
