package service

import (
	"app/internal"
	"errors"
	"fmt"
)
// en esta capa hago validaciones, me conecto con cosas esternas,
// y con el repository

// trabaja con la capa repository
type ProductDefault struct {

	rp internal.ProductRepository
}

// NewProductDefault crea un nuevo servicio de peliculas
func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

// Save guarda una pelicula en la base de datos y la retorna
func (m *ProductDefault) Save(product *internal.Product) (err error) {
	// valido la pelicula
	if (*product).Name == "" {
		return internal.ErrorProductNameRequired
	}
	if (*product).Quantity == 0 {
		return internal.ErrorProductQuantityRequired
	}
	if (*product).CodeValue == "" {
		return internal.ErrorProductCodeValueRequired
	}
	if (*product).Expiration == "" {
		return internal.ErrorProductExpirationRequired
	}
	if (*product).Price == 0.0 {
		return internal.ErrorProductPriceRequired
	}

	// save the product in the database
	err = m.rp.Save(product)
	if err != nil {
		switch err {
		case internal.ErrProductAlreadyExists:
			err = fmt.Errorf("product %s already exists", (*product).Name)
		default:
			err = fmt.Errorf("error saving product %s: %w", (*product).Name, err)
		}
	}
	return
}