package main

import (
	"net/http"

	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/internal/domain/service"
	"github.com/elaurentium/listener-net/internal/infra/persistence/db"
	"github.com/elaurentium/listener-net/internal/infra/web"
	"github.com/elaurentium/listener-net/internal/infra/web/auth"
	"github.com/elaurentium/listener-net/internal/infra/web/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Logger.Println("Initializing")

	dbConn, err := db.NewDBConnection()
	if err != nil {
		cmd.Logger.Println(err)
	}

	defer dbConn.Close()

	userRepo := db.NewUserRepository(dbConn)
	authService := auth.NewAuthService()
	userService := service.NewUserService(userRepo, authService)
	userHandler := handler.NewUserHandler(userService)

	router := web.NewRouter(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		cmd.Logger.Println(err)
	}

	cmd.Logger.Println("Run {}", server)
}
