package repository

import "app/internal"

type MovieMap struct {

	db map[int]internal.Movie

	lastId int
}
// esta capa se encarga de conectarse con la base de datos 


// creo un metodo newMovieMao que crea un nuevo map de peliculas
func NewMovieMap( db map[int]internal.Movie, lastId int) *MovieMap {
	return &MovieMap{
		db: db,
		lastId: lastId,
	}
}

// Save guarda una pelicula en la base de datos es decir, implementa el metodo Save de la interfaz MovieRepository
// y tambien hago las validaciones necesarias
func (m *MovieMap) Save(movie *internal.Movie) (err error) {

	// valido que el titulo de la pelicula no exista
	for _, v := range (*m).db {
		if v.Title == movie.Title {
			return internal.ErrorMovieTitleExists
		}
	}

	// incremento el id de la pelicula
	(*m).lastId++
	(*movie).Id = (*m).lastId

	// guardo la pelicula en la base de datos
	(*m).db[(*movie).Id] = *movie

	return
}