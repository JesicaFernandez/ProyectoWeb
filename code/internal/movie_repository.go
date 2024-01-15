package internal

import "errors"

var (
	ErrorMovieTitleExists = errors.New("movie title already exists")
)

type MovieRepository interface {
	// Save guardar una pelicula en la base de datos
	Save(movie *Movie) (err error)
}

