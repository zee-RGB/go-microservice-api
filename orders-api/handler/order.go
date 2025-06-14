package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Size:   size,
		Offset: cursor,
	})
	if err != nil {
		fmt.Println("failed to find all", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	id, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderObj, err := o.Repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, order.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(orderObj)
	if err != nil {
		fmt.Println("failed to marshal order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(
			http.StatusBadRequest)
		return
	}
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		fmt.Println("failed to parse order id:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.OrderStatus != "pending" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now

	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.UpdateByID(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to update order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to encode order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an Order by ID")
}
