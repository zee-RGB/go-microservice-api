package handler

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zee-RGB/orders-api/model"
	"github.com/zee-RGB/orders-api/repository/order"
)

type Order struct {
	Repo *order.RedisRepo
}

// CRUD operations methods
func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	order := model.Order{
		OrderID:     uint(rand.Uint64()),
		CustomerID:  body.CustomerID,
		LineItems:   body.LineItems,
		OrderStatus: "pending",
		CreatedAt:   &now,
	}

	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all Orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an Order by ID")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an Order by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an Order by ID")
}
