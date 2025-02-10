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
var ctx context.Context

func load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	err := load(configPath)
	if err != nil {
		log.Fatalf("failed to load config \"%v\": %v", configPath, err)
	}
	pgDSN, ok := os.LookupEnv(dsnEnvName)
	if !ok {
		panic("DSN environment variable not set")
	}

	ctx = context.Background()
	con, err = pgx.Connect(ctx, pgDSN)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	defer closeConOrLog(ctx, con)
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

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
