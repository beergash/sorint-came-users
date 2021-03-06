package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"users-admin/app/handlers"
	"users-admin/app/logger"
	"users-admin/config"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	logger.Info.Println("Initialization Db configuration")
	logger.Info.Printf("DB username %s\n", config.DB.Username)
	logger.Info.Printf("DB password %s\n", config.DB.Password)
	logger.Info.Printf("DB Name %s\n", config.DB.Name)
	logger.Info.Printf("DB Dialect %s\n", config.DB.Dialect)
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host))
	if err != nil {
		fmt.Println(err)
	}
	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()

}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/users", a.GetAllUsers)
	a.Get("/users/{userid}", a.FindUserById)
	a.Post("/create-user", a.CreateUser)
	a.Delete("/users/{userid}", a.DeleteUser)
	a.Put("/users/{userid}", a.ModifyUser)
	a.Post("/users/search", a.SearchUsers)
}

/*
** Projects Handlers
 */
func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	handlers.GetAllUsers(a.DB, w, r)
}
func (a *App) FindUserById(w http.ResponseWriter, r *http.Request) {
	handlers.FindUserById(a.DB, w, r)
}
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	handlers.CreateUser(a.DB, w, r)
}
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	handlers.DeleteUser(a.DB, w, r)
}
func (a *App) ModifyUser(w http.ResponseWriter, r *http.Request) {
	handlers.ModifyUser(a.DB, w, r)
}
func (a *App) SearchUsers(w http.ResponseWriter, r *http.Request) {
	handlers.SearchUsersByFilters(a.DB, w, r)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
