package router

import (
	"fmt"
	"net/http"
	"rakkiz-backend/views"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)


func Runserver(db *gorm.DB) {
  r := newRouter()

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", views.RegisterView(db))
	})

  fmt.Println("Server runs on port 8000")
	http.ListenAndServe(":8000", r)
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger)

	return r
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}