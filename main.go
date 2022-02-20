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
	productRouter := chi.NewRouter()
	productRouter.Get("/products", ph.GetProducts)
	//Update method
	productRouter.Put("/products/{id:[0-9]+}", ph.UpdateProducts)
	//Add method
	productRouter.Post("/products", ph.AddProducts)
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
