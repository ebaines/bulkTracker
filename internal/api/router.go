package api

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRouter() *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group("/api")
	{
		dayApiRouter := apiRouter.Group("/day")
		{
			dayApiRouter.GET("/:id", GetDay)
			dayApiRouter.POST("", AddDay)
			dayApiRouter.PUT("/:id", AddDay)
			dayApiRouter.DELETE("/:id", DeleteDay)
		}
		daysApiRouter := apiRouter.Group("/days")
		{
			daysApiRouter.GET("", GetDays)
		}
	}



	return r
}
