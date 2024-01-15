package main

import (
	"app/cmd/handlers"
	"net/http"
	"github.com/go-chi/chi/v5"
	"fmt"
)

func main() {
	
	// agrego las dependencias
	db, err := handlers.LoadProducts("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// creo un controller de products
	hc := handlers.NewControllerProducts(db, len(db))

	// creo un router
	r := chi.NewRouter()

	// registro una ruta y un handler
	r.Get("/ping", hc.Ping())

	// agrego Route para agrupar rutas
	r.Route("/products", func(r chi.Router) {
		r.Get("/", hc.Get())
		r.Get("/{id}", hc.GetById())
		r.Get("/search", hc.ProductSearch())
		r.Post("/", hc.Save())
	})

	// paso la url de esta forma http://localhost:8080/users?id=1

	// inicio el servidor web en el puerto 8080
	http.ListenAndServe(":8080", r)

}