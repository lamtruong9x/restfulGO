package handlers

import (
	"Resful/data"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	productsList := data.GetProducts()
	err := productsList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal", http.StatusInternalServerError)
	}
	data.AddProduct(newProduct)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	p.l.Println("ID", id)
	if err != nil {
		http.Error(rw, "Invalid URI unable to convert to numer", http.StatusBadGateway)
	}
	newProduct := &data.Product{}
	err = newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal", http.StatusInternalServerError)
	}
	err = data.UpdateProducts(id, newProduct)
	if err != nil {
		http.Error(rw, "Unable to update data", http.StatusInternalServerError)
	}
}
