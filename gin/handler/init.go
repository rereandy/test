package handler

import (
	"github.com/LearnGin/handler/person"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.Engine) {
	r.Handle("GET", "/ping", PingHandler())
	r.Handle("POST", "/person/create", person.CreatePersonHandler())
}
