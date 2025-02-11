package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/internal/repository"
	"github.com/Kosfedev/auth/internal/repository/user"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
)

const (
	configPath     = ".env"
	baseURL        = "localhost:8081"
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

var userRep repository.UserRepository

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config \"%v\": %v", configPath, err)
	}

	ctx := context.Background()
	cnf, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %s", err.Error())
	}
	pgDSN := cnf.DSN()

	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		log.Fatalf("failed to establish connection to \"%v\": %v", pgDSN, err.Error())
	}
	defer func(ctx context.Context, con *pgx.Conn) {
		if err := con.Close(ctx); err != nil {
			log.Printf("Error while closing connection: %+v\n", err)
		}
	}(ctx, con)
	userRep = user.NewRepository(con)

	r := chi.NewRouter()
	r.Post(usersPostfix, createUserHandler)
	r.Get(userPostfix, getUserHandler)
	r.Put(userPostfix, updateUserHandler)
	r.Delete(userPostfix, deleteUserHandler)

	server := http.Server{
		Addr:         baseURL,
		Handler:      r,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
