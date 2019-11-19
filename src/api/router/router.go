package router

import (
	v1 "APS/src/api/v1"

	_ "APS/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(mw...)
	g.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// main router
	apiv1 := g.Group("/api/v1")
	user := apiv1.Group("/user")
	{
		user.GET("", v1.GetUsers)
		user.POST("/add", v1.AddUser)
		user.DELETE("/:user_name", v1.DeleteUser)
	}

	// 公开的一些接口
	return g
}
