package main

import (
	"log"
	"net/http"
	"os"

	"utnianos-api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	h := handlers.New()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "UTNianos API", "endpoints": ["/foros", "/foro/{slug}", "/tema/{slug}", "/search", "/usuario/{username}"]}`))
	})

	r.Get("/foros", h.GetForums)

	r.Route("/foro", func(r chi.Router) {
		r.Get("/{slug}", h.GetForum)
	})

	r.Route("/tema", func(r chi.Router) {
		r.Get("/{slug}", h.GetTopic)
	})

	r.Route("/usuario", func(r chi.Router) {
		r.Get("/{username}", h.GetUser)
	})

	r.Get("/search", h.Search)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not found"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
