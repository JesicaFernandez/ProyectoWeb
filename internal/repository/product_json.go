package repository

import (
	"app/internal"
	"encoding/json"
)

func NewProductJSON(st internal.ProductStorage, db map[int]internal.Product, lastId int) *ProductJSON {
	// default config
	if db == nil {
		db = make(map[int]internal.Product)
	}

	// return the product map
	return &ProductJSON{
		db:     db,
		lastId: lastId,
		st:     st,
	}
}

// ProductJSON is the struct for the product map
type ProductJSON struct {
	//db is the map of products
	//- key is the id of the product
	//- value is the product
	db map[int]internal.Product
	//lastId is the last id of the product
	lastId int
	// storage is the storage of the product
	st internal.ProductStorage
}

func (p *ProductJSON) Save(product *internal.Product) (err error) {
	// read data from file
	err = p.readFile(&(*p).db)
	if err != nil {
		return
	}

	// validate the product
	// code_value is unique
	if err = p.ValidateCodeValueExistance((*product).CodeValue); err != nil {
		return
	}

	// increment the last id
	(*p).lastId++
	// set the id of the product
	(*product).Id = (*p).lastId

	// store the product in the map
	(*p).db[(*product).Id] = *product
	// convert the map to json
	data, err := json.Marshal((*p).db)
	if err != nil {
		return
	}

	// save data to file
	err = (*p).st.Write(data)
	if err != nil {
		return
	}

	return
}

// GetAll gets all the products from the repository
func (p *ProductJSON) GetAll() (products []internal.Product, err error) {
	// read data from file
	err = p.readFile(&(*p).db)
	if err != nil {
		return
	}

	// get all the products from the map
	for _, p := range (*p).db {
		products = append(products, p)
	}
	return
}

// GetByID gets the product by id from the repository
func (p *ProductJSON) GetByID(id int) (product internal.Product, err error) {
	// read data from file
	err = p.readFile(&(*p).db)
	if err != nil {
		return
	}

	// get the product from the map
	product, ok := (*p).db[id]
	if !ok {
		err = internal.ErrProductIdNotFound
		return
	}
	return
}

// Update updates the product in the repository
func (p *ProductJSON) Update(product *internal.Product) (err error) {
	// read data from file
	err = p.readFile(&(*p).db)
	if err != nil {
		return
	}

	// update the product in the map
	prod, ok := (*p).db[(*product).Id]
	if !ok {
		err = internal.ErrProductIdNotFound
		return
	}
	(*p).db[prod.Id] = *product

	// save data to file
	err = p.writeFile(&(*p).db)
	if err != nil {
		return
	}
	return
}

// Delete deletes the product from the repository
func (p *ProductJSON) Delete(id int) (err error) {
	// read data from file
	err = p.readFile(&(*p).db)
	if err != nil {
		return
	}

	// validate if the id exists
	if _, ok := (*p).db[id]; !ok {
		err = internal.ErrProductIdNotFound
		return
	} else {
		// delete the product from the map
		delete((*p).db, id)
	}

	// save data to file
	err = p.writeFile(&(*p).db)
	if err != nil {
		return
	}
	return
}

// ValidateIdExistance validates if the id exists
func (p *ProductJSON) ValidateIdExistance(id int) (err error) {
	// validate the id
	_, ok := (*p).db[id]
	if !ok {
		err = internal.ErrProductIdNotFound
		return
	}
	return
}

// ValidateCodeValueExistance validates if the code value exists
func (p *ProductJSON) ValidateCodeValueExistance(codeValue string) (err error) {
	// validate the code value
	for _, p := range (*p).db {
		if p.CodeValue == codeValue {
			err = internal.ErrProductCodeAlreadyExists
			return
		}
	}
	return
}

// ReadFile reads the file
func (p *ProductJSON) readFile(db *map[int]internal.Product) (err error) {
	// read data from file
	data, err := (*p).st.Read()
	if err != nil {
		return
	}

	// deserialize json to map
	if len(data) > 0 {
		err = json.Unmarshal(data, &(*p).db)
		if err != nil {
			return
		}
	}
	return
}

// WriteFile writes the file
func (p *ProductJSON) writeFile(db *map[int]internal.Product) (err error) {
	// convert the map to json
	data, err := json.Marshal((*p).db)
	if err != nil {
		return
	}

	// save data to file
	err = p.st.Write(data)
	return
}