package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"app/internal/storage"
	"net/http"
	"app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// NewDefaultHTTP creates a new default HTTP
func NewDefaultHTTP(address, token string) *DefaultHTTP {
	// return the default HTTP
	return &DefaultHTTP{
		address: address,
		token:   token,
	}
}

// DefaultHTTP is the struct for the default HTTP
type DefaultHTTP struct {
	// address is the address of the HTTP
	address string
	token   string
}

// Run runs the HTTP server
func (d *DefaultHTTP) Run() (err error) {
	// initialize dependencies
	// create the product storage
	st := storage.NewProductJSON(nil, "")
	// create the file storage
	err = st.CreateFile()
	if err != nil {
		return err
	}
	// create the product repository
	//rp := repository.NewProductMap(nil, 0)
	rp := repository.NewProductJSON(st, nil, 0)
	// create the product service
	sv := service.NewProductDefault(rp)
	// create the product handler
	hd := handler.NewProductDefault(sv)
	// create the chi router
	rt := chi.NewRouter()

	// add the authenticator middleware
	rt.Use(middleware.NewAuthenticator(d.token).Auth, middleware.NewResponseLogger().Log)

	// register the routes
	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll())
		r.Get("/{id}", hd.GetByID())
		r.Post("/", hd.Create())
		r.Put("/{id}", hd.Update())
		r.Patch("/{id}", hd.UpdatePartial())
		r.Delete("/{id}", hd.Delete())
		r.Get("/consumer_price", hd.ConsumerPrice())
	})

	// run the HTTP server
	err = http.ListenAndServe(d.address, rt)
	return
}