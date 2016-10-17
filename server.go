package main

import (
	"knowledge-base/router"
	"knowledge-base/routes"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	routes := routes.CreateRoutes()
	router := router.NewHatchRouter(routes)
	loggerRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	http.ListenAndServe(":3810", loggerRouter)
}
