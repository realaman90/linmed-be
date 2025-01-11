package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
)

func (s *Server) AddCustomer(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	customer := models.Customer{}

	// Unmarshal the request body into the customer struct
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate customer
	if err := customer.Validate(); err != nil {
		s.Logger.Error("Failed to validate customer: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// save to db
	id, err := s.db.AddCustomer(ctx, customer)
	if err != nil {
		s.Logger.Error("Failed to save customer to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      id,
		"message": "Customer added successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetCustomer(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get the customer id from the request
	id := mux.Vars(r)["id"]

	// validate id

	if id == "" {
		s.Logger.Error("Customer id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// Get the customer from the db
	customer, err := s.db.GetCustomer(ctx, id)
	if err != nil {
		s.Logger.Error("Failed to get customer from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, customer)
}

func (s *Server) GetCustomers(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")

	// validate page and limit query parameters
	pageInt, limitInt := s.validatePageLimit(page, limit)

	// Get all customers from the db

	customers, total, err := s.db.GetCustomers(ctx, pageInt, limitInt)
	if err != nil {
		s.Logger.Error("Failed to get customers from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := paginatedResponse{
		Total: total,
		Data:  customers,
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the customer id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("Customer id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	customer := models.Customer{}

	// Unmarshal the request body into the customer struct
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate customer
	if err := customer.Validate(); err != nil {
		s.Logger.Error("Failed to validate customer: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// update customer
	if err := s.db.UpdateCustomer(ctx, customer); err != nil {
		s.Logger.Error("Failed to update customer: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	res := map[string]string{
		"id":      id,
		"message": "Customer updated successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}
