package main

import (
	"github.com/LearnGin/handler"
	"github.com/LearnGin/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin with default configs
	r := gin.Default()

	// append custom middle-wares
	middleware.RegisterMiddleware(r)
	// register custom routers
	handler.RegisterHandler(r)

	// run the engine
	r.Run()
}
