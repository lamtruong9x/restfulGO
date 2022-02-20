package handlers

import (
	"Resful/data"
	"context"
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
	newProduct := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&newProduct)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	p.l.Println("ID", id)
	if err != nil {
		http.Error(rw, "Invalid URI unable to convert to number", http.StatusBadGateway)
	}

	newProduct := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProducts(id, &newProduct)
	if err != nil {
		http.Error(rw, "Unable to update data", http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

//middleware
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
