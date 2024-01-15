package internal

import "errors"

var (
	ErrorProductTitleExists = errors.New("product title already exists")
)

type ProductRepository interface {
	// Save guardar una pelicula en la base de datos
	Save(product *Product) (err error)
}

