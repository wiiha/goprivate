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

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return nil
}
