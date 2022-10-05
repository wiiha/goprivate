package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/wiiha/goprivate/server"
	"github.com/wiiha/goprivate/snote"
)

func main() {

	noteSVC := snote.NewNoteServiceWithProductionRepo()

	backend := server.NewServer("localhost:8080", noteSVC)
	backend.ListenAndServe()

}
