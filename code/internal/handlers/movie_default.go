package handlers

import (
	"net/http"
	"code/internal"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Crear NewMovie donde va a guardar las peliculas y va a tener un mapa de peliculas
// y un id para cada pelicula
func NewDefaultMovie(movies map[int]internal.Movie, lastId int) *DefaultMovie {
	
	// retornar un puntero a DefaultMovie para que sea accesible desde otros paquetes
	return &DefaultMovie{
		// retornar un mapa de peliculas
		movies: movies,
		// retornar un id para cada pelicula
		lastId: lastId,
	}
}

// DefaultMovie is a struct that represents a movie
type DefaultMovie struct {
	// crear nuestras peliculas, con un mapa
	movies map[int]internal.Movie

	// crear un id para cada pelicula
	lastId int
}

// creo movieJson para devolver la pelicula en formato json
type MovieJson struct {
	Id int 				`json:"id"`
	Title string 		`json:"title"`
	Year int 			`json:"year"`
	Publisher string 	`json:"publisher"`
	Rating float64 		`json:"rating"`
}

// creo un struct para que me envien la peticion y la obtengo del body (el cuerpo de la peticion)
type BodyRequest struct {
	Title string 		`json:"title"`
	Year int 			`json:"year"`
	Publisher string	`json:"publisher"`
	Rating float64 		`json:"rating"`
}

// crear metodo save para guardar las peliculas
func (h *DefaultMovie) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// creo el token
		token := r.Header.Get("Authorization")
		// valido el token
		if token != "123456" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}

		// validar request 
		
		//entonces agregamos un byte que guarde en memoria
		/*bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error reading the request body"))
			return
		}
		// valido la existencia del body
		// creo un mapa de tipo string y de tipo any
		var mp map[string]any

		if err := json.Unmarshal(bytes, &mp); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error falta el campo title, year, rating o publisher"))
			return
		}

		if err := validateExistence(mp, "title", "year", "rating", "publisher"); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error falta el campo title, year, rating o publisher"))
			return
		}*/

		// request ----------------------------------------------------
		// me envian la peticion y la obtengo del body (el cuerpo de la peticion)
		// creo una variable body de tipo BodyRequest
		var body BodyRequest
		// decodifico el body y lo guardo en la variable body
		// el decode lo que hace es que si coincide con el campo va a pasar el tipo de dato que esta en la estructura
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error decoding the request body"))
			return
		}

		// process ----------------------------------------------------
		// creo una pelicula
		// primero hay que serializar la pelicula

		(*h).lastId++
		id := (*h).lastId
		movie := internal.Movie{
			Id: id,
			Title: body.Title,
			Year: body.Year,
			Publisher: body.Publisher,
			Rating: body.Rating,
		}

		// agrego las validaciones de los campos de la pelicula
		if err := Validate(&movie); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error validating the movie"+err.Error()))
			return
		}

		// otro tipo de validación (si la pelicula ya existe)
		for _, m := range (*h).movies {
			if m.Title == movie.Title {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("movie already exists"))
				return
			}
		}

		// decodifico el body y lo guardo en la variable mp y valido que exista cada campo
		// con esta validación nos aseguramos que el cliente nos pase las claves
		// no me sirve mucho este codigo porque cuando llega ala otra validación la encuentra vacia
		/*if err := validateExistence(mp, "title", "year", "rating", "published"); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error validating the request body" + err.Error()))
			return
		}*/

		// response ----------------------------------------------------
		// serializo movieJson
		(*h).movies[movie.Id] = movie
		data := MovieJson{
			Id: movie.Id,
			Title: movie.Title,
			Year: movie.Year,
			Publisher: movie.Publisher,
			Rating: movie.Rating,
		}

		// seteo el header ahora que ya tengo la pelicula
		w.Header().Set("Content-Type", "application/json")
		// seteo el status code
		w.WriteHeader(http.StatusCreated)
		// devuelvo la pelicula en formato json con un mensaje de que se creo la pelicula
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Movie created successfully",
			"data": data,
		})

	}
}

// vamos a crear una funcion de validacion para validar los campos de la pelicula
func Validate(movie *internal.Movie) error {

	if movie.Title == "" {
		return errors.New("Title is required")
	}
	if movie.Year == 0 {
		return errors.New("Year is required")
	}
	if movie.Publisher == "" {
		return errors.New("Publisher is required")
	}
	if len(movie.Title) < 3 || len(movie.Title) > 100 {
		return errors.New("Title must be between 3 and 100 characters")
	}
	if movie.Year < 1900 || movie.Year > 2024 {
		return errors.New("Year must be between 1900 and 2021")
	}
	return nil
}

// validar la existencia de la pelicula
func validateExistence(mp map[string]any, keys ...string) error {
	
	// recorro el mapa de peliculas
	for _, key := range keys {
		// si la pelicula no existe
		if _, ok := mp[key]; !ok {
			return fmt.Errorf("key %s does no exist", key)
		}
	}
	return nil
}