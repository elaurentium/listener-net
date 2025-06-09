package main

import (
	"net/http"

	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/cmd/sub"
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
		cmd.Logger.Fatalf("Failed to connect to database: %v", err)
		return
	}

	defer dbConn.Close()

	cmd.Logger.Println("Database connected successfully")

	userRepo := db.NewUserRepository(dbConn)
	authService := auth.NewAuthService()
	userService := service.NewUserService(userRepo, authService)
	go sub.Interfaces(userService)
	userHandler := handler.NewUserHandler(userService)

	router := web.NewRouter(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	cmd.Logger.Printf("Server started on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		cmd.Logger.Fatalf("Server failed to start: %v", err)
	}
}
