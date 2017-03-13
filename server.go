package main

import (
	"knowledge-base/router"
	"knowledge-base/routes"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"knowledge-base/conf"
)

func main() {
	routes := routes.CreateRoutes()
	router := router.NewHatchRouter(routes)
	loggerRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	port := conf.Config.PORT
	if port == "" {
		port = "3810"
	}
	port = ":" + port

	http.ListenAndServe(port, loggerRouter)
}
