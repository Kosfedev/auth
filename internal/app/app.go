package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	desc "github.com/Kosfedev/auth/pkg/user_v1/gRPC"
	gRPCServer "github.com/Kosfedev/auth/pkg/user_v1/gRPC/server"
	"github.com/Kosfedev/auth/pkg/user_v1/http/handlers"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	configPath     = ".env"
	pathname       = "localhost"
	httpPort       = 8081
	grpcPort       = 8082
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

// App is...
type App struct {
	serviceProvider *serviceProvider
	router          *chi.Mux
}

// NewApp is...
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run is...
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var err error
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = a.runHTTPServer()
	}()

	go func() {
		defer wg.Done()
		err = a.runGRPCServer()
	}()

	wg.Wait()

	return err
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	// TODO: replace with console input
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	a.serviceProvider.UserImpl(ctx)

	r := chi.NewRouter()
	r.Post(usersPostfix, func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUserHandler(w, r, *a.serviceProvider.userImpl)
	})
	r.Get(userPostfix, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserHandler(w, r, *a.serviceProvider.userImpl)
	})
	r.Patch(userPostfix, func(w http.ResponseWriter, r *http.Request) {
		handlers.PutUserHandler(w, r, *a.serviceProvider.userImpl)
	})
	r.Delete(userPostfix, func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUserHandler(w, r, *a.serviceProvider.userImpl)
	})
	a.router = r

	return nil
}

func (a *App) runHTTPServer() error {
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", pathname, httpPort),
		Handler:      a.router,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	log.Printf("http server listening on :%d\n", httpPort)
	return server.ListenAndServe()
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", pathname, grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &gRPCServer.Server{})

	log.Printf("gRPC server listening on :%d\n", grpcPort)
	return s.Serve(lis)
}
