package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
registerRoutesToRouter handles attaching all
routes relevant for the application.
*/
func registerRoutesToRouter(router *gin.Engine) error {

	/*
		The landing page must be specified as
		a file in order to not create a
		`/*filepath` wildcard route that would
		consume all other routes.
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
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return nil
}
