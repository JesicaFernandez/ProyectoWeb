package internal

import (
	"errors"
)

var (
	ErrorMovieTitleRequired = errors.New("product title is required")
	ErrorMovieYearRequired = errors.New("product year is required")
	ErrorMoviePublisherRequired = errors.New("product publisher is required")
	ErrorMovieRatingRequired = errors.New("product rating is required")
	ErrMovieAlreadyExists = errors.New("product already exists")
)
var (
	ErrFieldRequired = errors.New("field required")
	ErrFieldQuality = errors.New("field quality")
)
// ProductService is an interface that represents a product service
type ProductService interface {
	// Save guardar un product en la base de datos
	Save(product *Product) (err error)
}
