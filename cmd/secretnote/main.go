package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/wiiha/goprivate/server"
)

func main() {

	backend := server.NewServer("localhost:8080")
	backend.ListenAndServe()

}
