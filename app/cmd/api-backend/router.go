package server

import (
	"net/http"

	"github.com/RodrigoMS/App_Go/cmd/api-backend/controllers"
	"github.com/RodrigoMS/App_Go/cmd/api-backend/views"
)

func routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /user", controllers.GetUser)
	router.HandleFunc("POST /user", controllers.PostUser)
	router.HandleFunc("PUT /user", controllers.PutUser)
	router.HandleFunc("DELETE /user", controllers.DeleteUser)

	router.HandleFunc("GET /product", controllers.GetProduct)
	router.HandleFunc("POST /product", controllers.PostProduct)
	router.HandleFunc("PUT /product", controllers.PutProduct)
	router.HandleFunc("DELETE /product", controllers.DeleteProduct)

	router.HandleFunc("/", views.HandleNotFound)

	return router
}
