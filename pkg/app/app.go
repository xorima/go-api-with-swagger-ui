package app

import (
	"context"
	_ "demo-api/docs"
	apifun "demo-api/pkg/app/api"
	"demo-api/pkg/utils/obervability"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Router interface {
	chi.Router
}

type App struct {
	router Router
}

func NewApp() *App {
	app := &App{
		router: chi.NewRouter(),
	}
	app.router.Handle("/metrics", promhttp.Handler())
	app.RegisterRoutes(app.router)

	return app

}

func (a *App) RegisterRoutes(router Router) {
	uh := &UserHandler{}
	uh.RegisterRoutes(router)
}

func (a *App) Run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := obervability.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      a.router,
	}

	// Start HTTP server.
	srvErr := make(chan error, 1)
	go func() {
		fmt.Printf("Starting Server on %s", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

type RouteHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(r Router)
}

// Response - Response struct
//
//	@name	Response
type Response struct {
	Status  int
	Message string
}

// User - User struct
//
//	@name	User
type User struct {
	Name  string
	Email string
}

// UserResponse - UserResponse struct
//
//	@name	UserResponse
type UserResponse struct {
	Response
	User User
}

type UserHandler struct {
}

// Get - Returns all the users
//
//	@Summary		This API can be used for User stuff.
//	@Description	Tells if the chi-swagger APIs are working or not.
//	@Tags			User Info
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	app.UserResponse	"Successful Response"
//	@Failure		404	{object}	app.Response		"Failure Response"
//	@Router			/users [get]
func (uh *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user get"))

}
func (uh *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user post"))
}
func (uh *UserHandler) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user put"))
}
func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user delete"))
}

func (uh *UserHandler) newHTTPHandler() http.Handler {
	mux := chi.NewMux()
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}
	handleFunc("/users", uh.Get)
	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func (uh *UserHandler) RegisterRoutes(r Router) {
	handleFunc := func(method, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))

		r.Method(method, pattern, handler)
	}
	handleFunc("GET", "/dice/roll", apifun.RollDice)

	handleFunc("GET", "/users", uh.Get)
	handleFunc("PUT", "/users", uh.Put)
	handleFunc("POST", "/users", uh.Post)
	handleFunc("DELETE", "/users", uh.Delete)

	swagHandler := otelhttp.WithRouteTag("/swagger", http.HandlerFunc(httpSwagger.WrapHandler))
	r.Mount("/swagger", swagHandler)

}
