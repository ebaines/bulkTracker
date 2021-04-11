package api

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRouter() *gin.Engine{
	r := gin.Default()

	apiRouter := r.Group("/api")
	{
		apiRouter.GET("/day/:id", GetDay)
		apiRouter.POST("/day", AddDay)
		apiRouter.PUT("/day/:id", AddDay)
		apiRouter.DELETE("/day/:id", DeleteDay)
	}

	return r
}
