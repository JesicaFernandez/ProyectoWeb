package handlers

import (
	"app/internal"
	"fmt"
	"encoding/json"
	"net/http"
)

// solo voy a tener por parametro el servicio de peliculas
func NewDefaultProduct(sv internal.ProductService) *DefaultProduct {
	// creo un nuevo DefaultProduct y le paso el servicio de peliculas
	return &DefaultProduct{
		sv: sv,
	}
}

// DefaultProduct is a struct that represents a product
type DefaultProduct struct {	
	// creo un campo de tipo ProductService para poder usar los metodos de la interfaz ProductService
	sv internal.ProductService
}

// creo productJson para devolver la pelicula en formato json
type ProductJson struct {
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
func (h *DefaultProduct) Save() http.HandlerFunc {
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
		product := internal.Product{
			Title: body.Title,
			Year: body.Year,
			Publisher: body.Publisher,
			Rating: body.Rating,
		}
	
		// para guardar la pelicula llamo al servicio de peliculas
		if err := h.sv.Save(&product); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error saving the product"))
			return
		}

		// response ----------------------------------------------------
		// serializo productJson
		data := ProductJson{
			Id: product.Id,
			Title: product.Title,
			Year: product.Year,
			Publisher: product.Publisher,
			Rating: product.Rating,
		}

		// seteo el header ahora que ya tengo la pelicula
		w.Header().Set("Content-Type", "application/json")
		// seteo el status code
		w.WriteHeader(http.StatusCreated)
		// devuelvo la pelicula en formato json con un mensaje de que se creo la pelicula
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Product created successfully",
			"data": data,
		})

	}
}

func loadProducts(filename string) []Product {
	
	filePath := "./docs/json/" + filename
	// Abre el archivo
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return nil
	}
	// Cierra el archivo cuando termine la función
	defer file.Close()

	// Lee el archivo y lo guarda en un slice de bytes
	// ioutil.ReadAll lee el archivo y lo guarda en un slice de bytes
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error leyendo el archivo:", err)
		return nil
	}

	// Creo una variable que es un slice de Product
	var products []Product
	// json.Unmarshal decodifica el json y lo guarda en la variable products
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		fmt.Println("Error decodificando el archivo JSON:", err)
		return nil
	}

	// retorno el slice de Product
	return products
}

func NewProduct() *MyProduct {
	return &MyProduct{
		Products: loadProducts("products.json"),
	}
}

/*Crear una ruta /products/:id que nos devuelva un producto por su id.*/

func (h *MyProduct) GetById(w http.ResponseWriter, r *http.Request) {
	
	// obtengo el id de la url
	id := chi.URLParam(r, "id")
	
	// busco el producto por el id
	var product *Product
    for _, p := range h.Products {
        if strconv.Itoa(p.Id) == id {
            product = &p
            break
        }
    }

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Product not found"))
		return
	}

	// Header para indicar que la respuesta es un json
	w.Header().Set("Content-Type", "application/json")
	// Encode el slice de Product a json y lo escribe en la respuesta
	p, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(p)
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