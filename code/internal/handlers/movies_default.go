package handlers

import (
	"app/internal"
	"fmt"
	"encoding/json"
	"net/http"
)

// solo voy a tener por parametro el servicio de peliculas
func NewDefaultMovie(sv internal.MovieService) *DefaultMovie {
	// creo un nuevo DefaultMovie y le paso el servicio de peliculas
	return &DefaultMovie{
		sv: sv,
	}
}

// DefaultMovie is a struct that represents a movie
type DefaultMovie struct {	
	// creo un campo de tipo MovieService para poder usar los metodos de la interfaz MovieService
	sv internal.MovieService
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

		// request
		// paso el json a una estructura
		var body BodyRequest
		// decodifico el body y lo guardo en la variable body
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error decoding the request body"))
			return
		}

		// process ----------------------------------------------------
		// creo una pelicula
		// primero hay que serializar la pelicula
		movie := internal.Movie{
			Title: body.Title,
			Year: body.Year,
			Publisher: body.Publisher,
			Rating: body.Rating,
		}
	
		// para guardar la pelicula llamo al servicio de peliculas
		if err := h.sv.Save(&movie); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error saving the movie"))
			return
		}

		// response ----------------------------------------------------
		// serializo movieJson
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