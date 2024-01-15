package service

import (
	"app/internal"
	"fmt"
	"time"
)

// NewProductDefault creates a new instance of ProductDefault
func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

// ProductDefault is the struct for the product default
type ProductDefault struct {
	// rp is the product repository
	rp internal.ProductRepository
}

// Save saves the product in the repository
func (p *ProductDefault) Save(product *internal.Product) (err error) {
	// business rules
	//- validate required fields
	if err = p.ValidateKeyContent(product); err != nil {
		return
	}

	// save the product in the repository
	if err = (*p).rp.Save(product); err != nil {
		switch err {
		case internal.ErrProductCodeAlreadyExists:
			err = fmt.Errorf("%w: code_value", internal.ErrProductAlreadyExists)
		}
		return
	}
	return
}

// GetAll gets all the products from the repository
func (p *ProductDefault) GetAll() (products []internal.Product, err error) {
	// get all the products from the repository
	products, err = (*p).rp.GetAll()
	return
}

// GetByID gets the product by id from the repository
func (p *ProductDefault) GetByID(id int) (product internal.Product, err error) {
	// get the product by id from the repository
	product, err = (*p).rp.GetByID(id)
	if err != nil {
		switch err {
		case internal.ErrProductIdNotFound:
			err = fmt.Errorf("%w: code_value", internal.ErrProductNotFound)
		}
		return
	}
	return
}

// Update updates the product in the repository
func (p *ProductDefault) Update(product *internal.Product) (err error) {
	// business rules
	//- validate required fields
	if err = p.ValidateKeyContent(product); err != nil {
		return
	}

	// update the product in the repository
	if err = (*p).rp.Update(product); err != nil {
		switch err {
		case internal.ErrProductIdNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}

// Delete deletes the product from the repository
func (p *ProductDefault) Delete(id int) (err error) {
	// delete the product from the repository
	if err = (*p).rp.Delete(id); err != nil {
		switch err {
		case internal.ErrProductIdNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}

// ConsumerPrice calculates the consumer price of the list of products
func (p *ProductDefault) ConsumerPrice(productIds *[]int) (products map[int]internal.Product, totalPrice float64, err error) {
	// totalQuantity is the total quantity of products
	var totalQuantity int

	// validate if the productIds is empty
	// - if is empty return the total price of the total quantity of products in the repository
	// - else return the total price of the productIds
	if len(*productIds) == 0 {
		// get all the products from the repository
		products, err := (*p).rp.GetAll()
		if err != nil {
			return nil, 0, err
		}
		// calculate the total price of the products
		for _, product := range products {
			// validate if the product is published
			if ValidateProductDisponibility(&product, 0) != nil {
				totalPrice += product.Price
				totalQuantity++
			}
		}
	} else {

		// process the ids for count the products requested
		productIdsCount := make(map[int]int)
		for _, productId := range *productIds {
			productIdsCount[productId]++
		}

		// calculate the total price of the productIdsMap
		for productId, quantity := range productIdsCount {
			// get the product by id from the repository
			product, err := (*p).rp.GetByID(productId)
			if err != nil {
				return nil, 0, err
			}

			// validate the product disponibility
			if err = ValidateProductDisponibility(&product, quantity); err != nil {
				return nil, 0, err
			}

			// calculate the total price
			totalPrice += product.Price * float64(quantity)
			// increase the total quantity
			totalQuantity += quantity
		}
	}

	// calculate tax on the total price
	switch {
	case totalQuantity < 10:
		totalPrice *= 1.21
	case totalQuantity < 20:
		totalPrice *= 1.17
	case totalQuantity > 20:
		totalPrice *= 1.15
	}
	
	return
}

// ValidateKeyContent validates the product fields with business rules
func (p *ProductDefault) ValidateKeyContent(product *internal.Product) (err error) {
	// validate the product
	//- name is not empty
	if (*product).Name == "" {
		err = fmt.Errorf("%w: name", internal.ErrFieldRequired)
		return
	}
	//- quantity is not negative
	if (*product).Quantity < 0 {
		err = fmt.Errorf("%w: quantity", internal.ErrInvalidField)
		return
	}
	//- code_value is not empty
	if (*product).CodeValue == "" {
		err = fmt.Errorf("%w: code_value", internal.ErrFieldRequired)
		return
	}
	//- expiration is not empty and is a valid date with format dd-mm-yyyy
	if !ValidateDateExpiration((*product).Expiration) {
		err = fmt.Errorf("%w: expiration", internal.ErrInvalidField)
		return
	}
	//- price is not negative or zero
	if (*product).Price <= 0 {
		err = fmt.Errorf("%w: price", internal.ErrInvalidField)
		return
	}
	return
}

// ValidateDateExpiration validates the expiration date with format dd/mm/yyyy
func ValidateDateExpiration(s string) (ok bool) {
	// parse the date
	date, err := time.Parse("02/01/2006", s)
	// validate the date
	if err != nil || date.Before(time.Now()) {
		return false
	}
	return true
}

// ValidateProductDisponibility validates the product with business rules
func ValidateProductDisponibility(product *internal.Product, quantity int) (err error) {
	// validate the product
	//- quantity is not negative and is less than the stock
	if quantity > (*product).Quantity {
		err = fmt.Errorf("%w: quantity", internal.ErrProductQuantityNotAvailable)
		return
	}
	if !(*product).IsPublished {
		err = fmt.Errorf("%w: is_published", internal.ErrProductNotPublished)
		return
	}
	return
}