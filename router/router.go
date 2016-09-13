package router

import (
	"knowledge-base/routes"

	"github.com/gorilla/mux"
)

func NewHatchRouter(routes routes.Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
