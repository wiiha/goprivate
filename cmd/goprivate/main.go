package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	flag "github.com/spf13/pflag"
	"github.com/wiiha/goprivate/server"
	"github.com/wiiha/goprivate/snote"
)

func main() {

	var hostAddress string
	var hostPort string
	var ginDebugMode bool
	var help bool

	flag.StringVarP(&hostAddress, "address", "a", "localhost", "address for host")
	flag.StringVarP(&hostPort, "port", "p", "8080", "port to serve on")
	flag.BoolVarP(&ginDebugMode, "debug", "d", false, "run gin (api framework) in debug mode")
	flag.BoolVarP(&help, "help", "h", false, "display this help message")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	gin.SetMode(gin.ReleaseMode)

	if ginDebugMode {
		gin.SetMode(gin.DebugMode)
	}

	noteSVC := snote.NewNoteServiceWithProductionRepo()

	addressAndPort := hostAddress + ":" + hostPort
	backend := server.NewServer(addressAndPort, noteSVC)
	log.Printf("success: serving on %s\n", addressAndPort)
	backend.ListenAndServe()

}
