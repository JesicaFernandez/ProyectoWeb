package main

import (
	"app/code/cmd/handlers"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	
	// creo un router
	r := chi.NewRouter()

	// creo un handler de productos
	h := handler.NewProduct()

	// registro una ruta y un handler
	r.Get("/ping", h.Ping())

	// agrego Route para agrupar rutas
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.GetAll)
		r.Get("/{id}", h.GetById)
		r.Get("/search", h.ProductSearch)
	})

	// paso la url de esta forma http://localhost:8080/users?id=1

	// inicio el servidor web en el puerto 8080
	http.ListenAndServe(":8080", r)

}