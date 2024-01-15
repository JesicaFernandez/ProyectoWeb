package handlers

import (
	"encoding/json"
	"fmt"
	"time"
	"errors"
	"regexp"
	"os"
)

type Product struct {
	Id          int    
	Name        string  
	Quantity    int     
	CodeValue   string  
	IsPublished bool    
	Expiration  time.Time  
	Price       float64 
}

type ProductAttributesJSON struct {
	Name 	  	string 		`json:"name"`
	Quantity  	int 		`json:"quantity"`
	CodeValue 	string 		`json:"code_value"`
	IsPublished bool 		`json:"is_published"`
	Expiration 	time.Time 	`json:"expiration"`
	Price 		float64		`json:"price"`
}

var (
	ErrValidateRequiredField = errors.New("validate required field")
	ErrorValidateQualityField = errors.New("validate quantity field")
)

func Validate(p *Product) (err error) {
	if p.Name == "" {
		err = fmt.Errorf("%w: name", ErrValidateRequiredField)
		return
	}
	if p.CodeValue == "" {
		err = fmt.Errorf("%w: name", ErrValidateRequiredField)
		return
	}
	if p.Expiration.IsZero() {
		err = fmt.Errorf("%w: name", ErrValidateRequiredField)
		return
	}
	
	if p.Quantity < 0 {
		err = fmt.Errorf("%w: quantity", ErrorValidateQualityField)
		return
	}
	if p.Price < 0 {
		err = fmt.Errorf("%w: price", ErrorValidateQualityField)
		return
	}
	if p.Expiration.Before(time.Now()) {
		err = fmt.Errorf("%w: expiration", ErrorValidateQualityField)
		return
	}
	rx := regexp.MustCompile(`^[A.Z]{3}-[0-9]{3}$`)
	if !rx.MatchString(p.CodeValue) {
		err = fmt.Errorf("%w: code_value", ErrorValidateQualityField)
		return
	}
	return
}

func LoadProducts(filename string) (p map[int]*Product, err error) {
	
	filePath := "./docs/json/" + filename
	// Abre el archivo
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return
	}
	// Cierra el archivo cuando termine la función
	defer file.Close()

	// Lee el archivo e inicializa un slice de bytes
	products := make(map[int]*ProductAttributesJSON)
	// Decodifica el json y lo guarda en la variable products
	err = json.NewDecoder(file).Decode(&products)
	if err != nil {
		fmt.Println("Error decodificando el archivo:", err)
		return
	}

	// crea un map de productos
	p = make(map[int]*Product)
	// recorre el map de productos y los agrega al map de productos
	for k, v := range products {
		p[k] = &Product{
			Name: 			v.Name,
			Quantity: 		v.Quantity,
			CodeValue: 		v.CodeValue,
			IsPublished: 	v.IsPublished,
			Expiration: 	v.Expiration,
			Price: 			v.Price,
		}
	}
	return
}