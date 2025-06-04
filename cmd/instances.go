package cmd

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/service"
	"github.com/elaurentium/listener-net/internal/infra/persistence/db"
	"github.com/elaurentium/listener-net/internal/infra/web"
	"github.com/elaurentium/listener-net/internal/infra/web/auth"
	"github.com/elaurentium/listener-net/internal/infra/web/handler"
	"github.com/joho/godotenv"
)

func Instance() {
	godotenv.Load()
	Logger.Println("Initializing")

	pool, err := db.NewDBConnection()
	if err != nil {
		Logger.Println(err)
		return
	}

	defer pool.Close()

	userRepo := db.NewUserRepository(pool)
	authService := auth.NewAuthService()
	userService := service.NewUserService(userRepo, authService)
	userHandler := handler.NewUserHandler(userService)

	router := web.NewRouter(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		Logger.Println(err)
	}
}
