package handlers

import (
	"Resful/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.AddProducts(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		rexp := regexp.MustCompile(`/([0-9]+)`)
		g := rexp.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.UpdateProducts(id, rw, r)
		return
	}
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

func (p *Products) UpdateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal", http.StatusInternalServerError)
	}
	err = data.UpdateProducts(id, newProduct)
	if err != nil {
		http.Error(rw, "Unable to update data", http.StatusInternalServerError)
	}
}
