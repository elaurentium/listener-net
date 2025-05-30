package web

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/infra/web"
	"github.com/elaurentium/listener-net/internal/domain/infra/web/handler"
)


func NewRouter(userHandler *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// web.CorsMiddleware(mux)

	// mux.HandleFunc("/users/{id}", userHandler)

	return mux
}