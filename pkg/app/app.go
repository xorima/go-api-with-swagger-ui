package app

import (
	_ "demo-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
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

func (a *App) Run() error {
	return http.ListenAndServe(":3000", a.router)
}

type RouteHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(r Router)
}

// Response - Response struct
//	@name	Response
type Response struct {
	Status  int
	Message string
}

// User - User struct
//	@name	User
type User struct {
	Name  string
	Email string
}

// UserResponse - UserResponse struct
//	@name	UserResponse
type UserResponse struct {
	Response
	User User
}

type UserHandler struct {
}

// Get - Returns all the users
//	@Summary		This API can be used for User stuff.
//	@Description	Tells if the chi-swagger APIs are working or not.
//	@Tags			User Info
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	app.UserResponse	"Successful Response"
//	@Failure		404	{object}	app.Response		"Failure Response"
//	@Router			/user [get]
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

func (uh *UserHandler) RegisterRoutes(r Router) {
	r.Get("/users", uh.Get)
	r.Post("/users", uh.Post)
	r.Put("/users", uh.Put)
	r.Delete("/users", uh.Delete)
	r.Mount("/swagger", httpSwagger.WrapHandler)
}
