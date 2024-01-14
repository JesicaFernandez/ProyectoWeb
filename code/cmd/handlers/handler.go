package handler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"github.com/go-chi/chi/v5"
	"os"
)

/*Nombre - Tipo de dato JSON - Tipo de dato GO - Descripción | Ejemplo

id - number - int - Identificador en conjunto de datos | 15

name - string - string - Nombre caracterizado | Cheese - St. Andre

quantity - number - int - Cantidad almacenada | 60

code_value - string - string - Código alfanumérico característico | S73191A

is_published - boolean - bool - El producto se encuentra publicado o no |  True

expiration - string - string - Fecha de vencimiento | 12/04/2022

price - number - float64 - Precio del producto | 50.15*/

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type MyProduct struct {
	Products []Product

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

func (h *MyProduct) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	}
}

// GetAll devuelve todos los productos
func (h *MyProduct) GetAll(w http.ResponseWriter, r *http.Request) {
	// Header para indicar que la respuesta es un json
	w.Header().Set("Content-Type", "application/json")
	// Encode el slice de Product a json y lo escribe en la respuesta
	products, err := json.Marshal(h.Products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(products)

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


//Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.

func (h *MyProduct) ProductSearch(w http.ResponseWriter, r *http.Request) {

	// busco el parametro priceGt en la url
	price := r.URL.Query().Get("price")

	// convierto priceGt a float64
	priceGt, err := strconv.ParseFloat(price, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}

	// busco el producto priceGt
	var productList []Product
	for _, p := range h.Products {
		if p.Price > priceGt {
			productList = append(productList, p)
		}
	}

	if productList == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Product not found"))
		return
	}

	// Header para indicar que la respuesta es un json
	w.Header().Set("Content-Type", "application/json")
	// Encode el slice de Product a json y lo escribe en la respuesta
	p, err := json.Marshal(productList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(p)


}