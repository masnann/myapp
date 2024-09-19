package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/username/myapp/handler"
	"github.com/username/myapp/repository"
	"github.com/username/myapp/service"
)

func main() {
	// Konfigurasi database
	dsn := "postgres://postgres:mkpmobile2024@localhost:5432/myapp_test?sslmode=disable"

	// Inisialisasi database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inisialisasi repository, service, dan handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Inisialisasi Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routing
	e.DELETE("/users/:id", userHandler.DeleteUser)

	// Jalankan server
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Shutting down the server: %v", err)
	}
}