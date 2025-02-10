package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	configPath     = ".env"
	dsnEnvName     = "PG_DSN"
	baseURL        = "localhost:8081"
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

var con *pgx.Conn
var pgDSN string

func load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func initDeps() {
	err := load(configPath)
	if err != nil {
		log.Fatalf("failed to load config \"%v\": %v", configPath, err)
	}
	pgDSN = os.Getenv(dsnEnvName)
	if len(pgDSN) == 0 {
		log.Fatalf("DSN environment variable not set")
	}
}

func main() {
	initDeps()
	var err error
	ctx := context.Background()
	con, err = pgx.Connect(ctx, pgDSN)
	if err != nil {
		log.Fatalf("failed to establish connection to \"%v\": %v", pgDSN, err.Error())
	}
	defer func(ctx context.Context, con *pgx.Conn) {
		if err := con.Close(ctx); err != nil {
			log.Printf("Error while closing connection: %+v\n", err)
		}
	}(ctx, con)

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
