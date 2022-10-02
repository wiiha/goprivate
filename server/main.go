package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	addr       string
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
func NewServer(addressAndPort string) *Server {
	router := gin.Default()

	err := registerRoutesToRouter(router)
	if err != nil {
		log.Fatalf("registering routes: %v", err)
	}

	s := &http.Server{
		Addr:         addressAndPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		addr:       addressAndPort,
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
