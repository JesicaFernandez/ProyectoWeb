package repository

import "app/internal"

type ProductMap struct {

	db map[int]internal.Product

	lastId int
}
// esta capa se encarga de conectarse con la base de datos 


// creo un metodo newProductMao que crea un nuevo map de peliculas
func NewProductMap( db map[int]internal.Product, lastId int) *ProductMap {
	return &ProductMap{
		db: db,
		lastId: lastId,
	}
}

// Save guarda una pelicula en la base de datos es decir, implementa el metodo Save de la interfaz ProductRepository
// y tambien hago las validaciones necesarias
func (m *ProductMap) Save(product *internal.Product) (err error) {

	// valido que el titulo de la pelicula no exista
	for _, v := range (*m).db {
		if v.Title == product.Title {
			return internal.ErrorProductTitleExists
		}
	}

	// incremento el id de la pelicula
	(*m).lastId++
	(*product).Id = (*m).lastId

	// guardo la pelicula en la base de datos
	(*m).db[(*product).Id] = *product

	return
}