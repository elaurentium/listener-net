package web

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/infra/web/handler"
)


func NewRouter(userHandler *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()


	return mux
}