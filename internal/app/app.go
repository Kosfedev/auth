package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/pkg/user_v1/http/handlers"
	"github.com/go-chi/chi"
)

const (
	configPath     = ".env"
	baseURL        = "localhost:8081"
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

type App struct {
	serviceProvider *serviceProvider
	router          *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
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
	r.Put(userPostfix, func(w http.ResponseWriter, r *http.Request) {
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
		Addr:         baseURL,
		Handler:      a.router,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	return server.ListenAndServe()
}
