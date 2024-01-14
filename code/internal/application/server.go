package application

import (
	"code/internal"
	"code/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// incializo todas las dependencia y ejecuto la aplicacion

// pedirle al usuario que ingrese el puerto y el host
type Server struct {
	address string
}

func NewServer(address string) *Server {
	// si el usuario no ingresa el puerto y el host, se va a ejecutar en el puerto 8080
	defaultAddress := ":8080"
	if address != "" {
		defaultAddress = address
	}
	
	return &Server{
		address: defaultAddress,
	}
}

// vamos a crear un metodo Run que va a ejecutar la aplicacion
func (s *Server) Run() error {

	// primero cargo las dependencias
	// creo un map de peliculas
	db := make(map[int]internal.Movie, 0)
	lastId := 0

	// creo un handler de peliculas
	hd := handlers.NewDefaultMovie(db, lastId)

	// creo un router
	rt := chi.NewRouter()

	// defino las rutas
	rt.Post("/movies", hd.Save())

	// ejecuto la aplicacion
	return http.ListenAndServe(s.address, rt)
}
