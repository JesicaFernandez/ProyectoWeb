package handlers

import (
	"time"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/go-chi/chi/v5"
)

type ControllerProducts struct {
	storage map[int]*Product
	lastId int
}

type ProductJSON struct {
	Id          int    		`json:"id"`
	Name        string 		`json:"name"`
	Quantity    int    		`json:"quantity"`
	CodeValue   string 		`json:"code_value"`
	IsPublished bool   		`json:"is_published"`
	Expiration  string 		`json:"expiration"`
	Price 		float64		`json:"price"`
}

type RequestBodyProductSave struct {
	Name        string 		`json:"name"`
	Quantity    int    		`json:"quantity"`
	CodeValue   string 		`json:"code_value"`
	IsPublished bool   		`json:"is_published"`
	Expiration  string 		`json:"expiration"`
	Price 		float64		`json:"price"`
}

type ResponseBodyProductSave struct {
	Message 	string 			`json:"message"`
	Data 		*ProductJSON 	`json:"data"`
}

type ResponseBodyProductGet struct {
	Message 	string 			`json:"message"`
	Data 		[]*ProductJSON 	`json:"data"`
}

type ResponseBodyProductGetById struct {
	Message 	string 			`json:"message"`
	Data 		*ProductJSON 	`json:"data"`
}

type ResponseBodyProductSearch struct {
	Message 	string 			`json:"message"`
	Data 		[]*ProductJSON 	`json:"data"`
}

func NewControllerProducts(storage map[int]*Product, lastId int) *ControllerProducts {
	return &ControllerProducts{
		storage: storage,
		lastId: lastId,
	}
}

func (c *ControllerProducts) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request -----------------------------------

		// process -----------------------------------

		// response ----------------------------------

		code := http.StatusOK
		// creo un body de respuesta que retorna un slice de productos
		body := ResponseBodyProductGet{
			Message: "success",
			Data:    make([]*ProductJSON, 0, len(c.storage)),
		}

		// recorro el storage y agrego los productos al body
		for k, v := range c.storage {
			// agrego el producto al body
			body.Data = append(body.Data, &ProductJSON{
				Id:          k,
				Name:        v.Name,
				Quantity:    v.Quantity,
				CodeValue:   v.CodeValue,
				IsPublished: v.IsPublished,
				Expiration:  v.Expiration.Format("2006-01-02"),
				Price:       v.Price,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

func (c *ControllerProducts) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		// request -----------------------------------
		// obtengo el id de la url, strconv.Atoi convierte un string a int
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductGetById{
				Message: "invalid id",
				Data:    nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}
		
		// process -----------------------------------
		product, ok := c.storage[id]
		if !ok {
			code := http.StatusNotFound
			body := ResponseBodyProductGetById{
				Message: "product not found",
				Data:    nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return 
		}

		// response ----------------------------------

		code := http.StatusOK
		// creo un body de respuesta que retorna un producto
		body := ResponseBodyProductGetById{
			Message: "success",
			Data:    &ProductJSON{
				Id: 			id,
				Name: 			product.Name,
				Quantity: 		product.Quantity,
				CodeValue: 		product.CodeValue,
				IsPublished: 	product.IsPublished,
				Expiration: 	product.Expiration.Format("2006-01-02"),
				Price: 			product.Price,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

func (h *ControllerProducts) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	}
}

//Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.

func (c *ControllerProducts) ProductSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		// request -----------------------------------
		// obtengo el precio de la url
		price := r.URL.Query().Get("price")
		// convierto price a float64
		priceGt, err := strconv.ParseFloat(price, 64)
		if err != nil {
			code := http.StatusBadRequest
			body:= ResponseBodyProductSearch{
				Message: 	"invalid price parameter",
				Data: 		nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		// process -----------------------------------
		// filtro los productos por price
		filterProducts := make([]*ProductJSON, 0)
		for k, v := range c.storage {
			if v.Price > priceGt {
				filterProducts = append(filterProducts, &ProductJSON{
					Id: 			k,
					Name: 			v.Name,
					Quantity: 		v.Quantity,
					CodeValue: 		v.CodeValue,
					IsPublished: 	v.IsPublished,
					Expiration: 	v.Expiration.Format("2006-01-02"),
					Price: 			v.Price,
				})
			}
		}
		code := http.StatusOK
		body := ResponseBodyProductSearch{
			Message: "success",
			Data:   filterProducts,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

func (c *ControllerProducts) Save() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		// request -----------------------------------
		var requestProduct RequestBodyProductSave
		
		//
		if err := json.NewDecoder(r.Body).Decode(&requestProduct); err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductSave{
				Message: 	"invalid request body",
				Data: 		nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		// process -----------------------------------
		exp, err := time.Parse("2006-01-02", requestProduct.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductSave{
				Message: 	"invalid date format. Must be yyy-mm-dd",
				Data: 		nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		product := &Product {
			Name: 			requestProduct.Name,
			Quantity:		requestProduct.Quantity,
			CodeValue:		requestProduct.CodeValue,
			IsPublished:	requestProduct.IsPublished,
			Expiration: 	exp,
			Price:			requestProduct.Price,
		}

		// validamos 
		if err := Validate(product); err != nil {
			code := http.StatusConflict
			body := ResponseBodyProductSave{
				Message: 	"invalid product" + err.Error(), 
				Data: 		nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}	

		// valido que el id sea unico

		/*if !IsUniqueID(product.Id, c.storage) {
			code := http.StatusConflict
			body := ResponseBodyProductSave{
				Message: "Product with ID already exists",
				Data:    nil,
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}*/

		c.lastId++
		c.storage[c.lastId] = product	

		// valido que el code_value sea unico por product
		/*if !IsUniqueCodeValue(product.CodeValue, c.storage) {
			code := http.StatusConflict
			body := ResponseBodyProductSave{
				Message: "Product with CodeValue already exists",
				Data:    nil,
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}*/
		// response ----------------------------------

		code := http.StatusCreated
		body := ResponseBodyProductSave{
			Message: 	"product created",
			Data:		&ProductJSON{
				Id: 			c.lastId,
				Name: 			product.Name,
				Quantity:		product.Quantity,
				CodeValue:		product.CodeValue,
				IsPublished:	product.IsPublished,
				Expiration:		product.Expiration.Format("2006-01-02"),
				Price:			product.Price,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}