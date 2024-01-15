package application

import (
	"app/internal"
	"app/internal/handlers"
	"app/internal/repository"
	"app/internal/service"
	"net/http"
	"github.com/go-chi/chi/v5"
)
type Default struct {
	addr string
}

func NewDefault(addr string) *Default {
	return &Default{
		addr: addr,
	}
}

func (d *Default) Run() (err error) {

	// primero inicializo las dependencias
	// creo el repositorio de peliculas
	rp := repository.NewMovieMap(make(map[int]internal.Movie), 0)

	// creo el servicio de peliculas
	sv := service.NewMovieDefault(rp)

	// creo el handler de peliculas
	hd := handlers.NewDefaultMovie(sv)

	rt := chi.NewRouter()

	// agrego rutas
	rt.Post("/product", hd.Save())
	rt.Get("/products/{id}", hd.GetById())

	err = http.ListenAndServe(d.addr, rt)
	return
}