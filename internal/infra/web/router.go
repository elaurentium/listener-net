package web

import (
	"github.com/elaurentium/listener-net/internal/infra/web/handler"
	"github.com/gin-gonic/gin"
)


func NewRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/user", userHandler.Register)
	router.GET("/user/:id", userHandler.GetUser)
	router.GET("/user/dispositive/:dispositive", userHandler.GetUserDispositive)

	return router
}