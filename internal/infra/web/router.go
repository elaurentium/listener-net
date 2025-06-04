package web

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/infra/web/handler"
)


func NewRouter(userHandler *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/user/:id", userHandler.GetUser)
	mux.HandleFunc("/user/dispositive/:dispositive", userHandler.GetUserDispositive)

	return mux
}