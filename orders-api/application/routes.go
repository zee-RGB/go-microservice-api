package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zee-RGB/orders-api/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRoutes)

	return router
}

func loadOrderRoutes(router chi.Router) {
	order.Handler := &handler.Order{}

	router.Post("/", order.Create)
	router.Get("/", order.List)
	router.Get("/{id}", order.GetByID)
	router.Put("/{id}", order.UpdateByID)
	router.Delete("/{id}", order.DeleteByID)
}
