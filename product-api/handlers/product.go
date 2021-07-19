package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/albukhary/learn-microservices/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		// Expect the ID from the URI
		rx := regexp.MustCompile(`/([0-9]+)`)

		g := rx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one ID")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one Capture Group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI. Unable to convert to Number")
			p.l.Println("Invalid ID")
		}

		p.l.Println("got id", id)

		p.updateProducts(id, rw, r)
		return
	}

	// catch all others

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod : %#v", prod)
	data.AddProduct(prod)
}

func (p Products) updateProducts (id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err2 := data.UpdateProduct(id, prod)
	if err2 == data.ErrProductNotFound {
		http.Error(rw, "Product Not found", http.StatusNotFound)
		return
	}

	if err2 != nil {
		http.Error(rw, "Product Not found", http.StatusInternalServerError)
		return
	}
}