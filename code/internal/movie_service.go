package internal

import (
	"errors"
)

var (
	ErrorMovieTitleRequired = errors.New("movie title is required")
	ErrorMovieYearRequired = errors.New("movie year is required")
	ErrorMoviePublisherRequired = errors.New("movie publisher is required")
	ErrorMovieRatingRequired = errors.New("movie rating is required")
	ErrMovieAlreadyExists = errors.New("movie already exists")
)
var (
	ErrFieldRequired = errors.New("field required")
	ErrFieldQuality = errors.New("field quality")
)
// MovieService is an interface that represents a movie service
type MovieService interface {
	// Save guardar una pelicula en la base de datos
	Save(movie *Movie) (err error)
}
