package server

import (
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
		API for backend
	*/
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return nil
}
