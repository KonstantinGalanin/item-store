package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/KonstantinGalanin/itemStore/internal/handlers"
	repository "github.com/KonstantinGalanin/itemStore/internal/repository/user"
	"github.com/KonstantinGalanin/itemStore/internal/router"
	"github.com/KonstantinGalanin/itemStore/internal/service"
	"github.com/KonstantinGalanin/itemStore/pkg/jwt"

	_ "github.com/lib/pq"
)

var (
	dbHost = os.Getenv("DATABASE_HOST")
	dbPort = os.Getenv("DATABASE_PORT")
	dbUser = os.Getenv("DATABASE_USER")
	dbPass = os.Getenv("DATABASE_PASSWORD")
	dbName = os.Getenv("DATABASE_NAME")

	serverPort = os.Getenv("SERVER_PORT")
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	userRepo := repository.NewUserPostgresRepo(db)
	userService := service.NewUserService(userRepo)

	jwtService := jwt.NewJwtService()
	userHandler := handlers.NewUserHandler(userService, jwtService)

	r := router.NewRouter(userHandler)
	err = http.ListenAndServe(":" + serverPort, r)
	if err != nil {
		panic(err)
	}
}
