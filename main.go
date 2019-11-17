package main

import (
	src "APS/src"
	router "APS/src/api/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


// @title APS API
// @version 1.0
// @description APS openAPI docs
// @contact.name API Support
// @contact.url http://tangmengqiu.github.io
// @contact.email sctmq@zju.edu.cn

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Create the Gin engine.
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	go src.CheckAt24()
	go src.CheckCommit()
	fmt.Println(http.ListenAndServe(":8080", g).Error())
	for {

	}
	return
}
