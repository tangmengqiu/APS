package main

import (
	src "APS/src"
	router "APS/src/api/router"
	"APS/tools/storage"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	cfg = flag.String("config", "", "APS config file path.")
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

	if err := src.InitConfig(*cfg); err != nil {
		panic(err)
	}
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)
	var err error
	src.MDataBase, err = storage.NewStorage()
	if err != nil {
		logrus.WithField("event", "new MDataBase")
		return
	}

	src.SyncMemoryToUsers()
	go src.CheckAt24(src.Check)
	go src.CheckCommit()
	fmt.Println(http.ListenAndServe(":8080", g).Error())
	for {

	}
}
