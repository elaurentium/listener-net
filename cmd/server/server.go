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
		cmd.Logger.Fatalf("Failed to connect to database: %v", err)
		return
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		cmd.Logger.Fatalf("Database ping failed: %v", err)
		return
	}
	cmd.Logger.Println("Database connected successfully")

	userRepo := db.NewUserRepository(dbConn)
	authService := auth.NewAuthService()
	userService := service.NewUserService(userRepo, authService)
	userHandler := handler.NewUserHandler(userService)

	router := web.NewRouter(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		cmd.Logger.Printf("Server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed  {
			cmd.Logger.Fatalf("Server failed to start: %v", err)
		}
	}()
}
