package main

import (
	"Resful/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	l := log.New(os.Stdout, "product-api-", log.LstdFlags)
	//create handlers
	ph := handlers.NewProducts(l)
	//Get product

	r.Get("/products", ph.GetProducts)

	r.Group(func(r chi.Router) {
		r.Use(ph.MiddlewareValidateProduct)
		//Update method
		r.Put("/products/{id:[0-9]+}", ph.UpdateProducts)
		//Add method
		r.Post("/products", ph.AddProducts)
	})
	//Mount productRouter to main router
	//r.Mount("/products", productRouter)

	sv := &http.Server{
		Addr:         ":9090",
		Handler:      r,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		IdleTimeout:  time.Second * 120,
	}
	go func() {
		err := sv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Graceful shutdown...", sig)
	tc, _ := context.WithTimeout(context.Background(), time.Second*30)
	sv.Shutdown(tc)
}
