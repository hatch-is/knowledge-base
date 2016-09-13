package main

import (
	"knowledge-base/router"
	"knowledge-base/routes"
	"net/http"
)

func main() {
	routes := routes.CreateRoutes()
	router := router.NewHatchRouter(routes)
	http.ListenAndServe(":3810", router)
}
