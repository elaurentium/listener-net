package web

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/infra/web/handler"
)


func NewRouter(userHandler *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.

	return mux
}