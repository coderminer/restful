package routes

import (
	"net/http"

	"github.com/coderminer/restful/controllers"
	"github.com/gorilla/mux"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

var routes []Route

func init() {
	register("GET", "/movies", controllers.AllMovies)
	register("GET", "/movies/{id}", controllers.FindMovie)
	register("POST", "/movies", controllers.CreateMovie)
	register("PUT", "/movies", controllers.UpdateMovie)
	register("DELETE", "/movies/{id}", controllers.DeleteMovie)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.Methods(route.Method).
			Path(route.Pattern).
			Handler(route.Handler)
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc) {
	routes = append(routes, Route{method, pattern, handler})
}
