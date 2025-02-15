package main

import (
	"database/sql"
	"net/http"

	"github.com/KonstantinGalanin/itemStore/internal/handlers"
	repository "github.com/KonstantinGalanin/itemStore/internal/repository/user"
	"github.com/KonstantinGalanin/itemStore/internal/router"
	"github.com/KonstantinGalanin/itemStore/internal/service"
	"github.com/KonstantinGalanin/itemStore/pkg/jwt"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=admin password=mypassword dbname=itemstore sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	userRepo := repository.NewUserPostgresRepo(db)
	userService := &service.UserService{
		UserRepo: userRepo,
	}

	jwtService := jwt.NewJwtService()
	userHandler := &handlers.UserHandler{
		UserService: userService,
		JwtService: jwtService,
	}

	r := router.NewRouter(userHandler)
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
