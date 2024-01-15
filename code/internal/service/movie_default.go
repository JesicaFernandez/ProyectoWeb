package service

import (
	"app/internal"
	"errors"
	"fmt"
)
// en esta capa hago validaciones, me conecto con cosas esternas,
// y con el repository

// trabaja con la capa repository
type MovieDefault struct {

	rp internal.MovieRepository
}

// NewMovieDefault crea un nuevo servicio de peliculas
func NewMovieDefault(rp internal.MovieRepository) *MovieDefault {
	return &MovieDefault{
		rp: rp,
	}
}

// Save guarda una pelicula en la base de datos y la retorna
func (m *MovieDefault) Save(movie *internal.Movie) (err error) {

	// valido la pelicula
	if (*movie).Title == "" {
		return internal.ErrorMovieTitleRequired
	}
	if (*movie).Year == 0 {
		return internal.ErrorMovieYearRequired
	}
	if (*movie).Publisher == "" {
		return internal.ErrorMoviePublisherRequired
	}
	if len((*movie).Title) < 3 || len((*movie).Title) > 100 {
		return errors.New("Title must be between 3 and 100 characters")
	}
	if (*movie).Year < 1900 || (*movie).Year > 2024 {
		return errors.New("Year must be between 1900 and 2021")
	}

	// save the movie in the database
	err = m.rp.Save(movie)
	if err != nil {
		switch err {
		case internal.ErrMovieAlreadyExists:
			err = fmt.Errorf("movie %s already exists", (*movie).Title)
		default:
			err = fmt.Errorf("error saving movie %s: %w", (*movie).Title, err)
		}
	}
	return
}