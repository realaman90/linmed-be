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
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate customer
	if err := customer.Validate(); err != nil {
		s.Logger.Error("Failed to validate customer: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// save to db
	if err := s.db.AddCustomer(ctx, customer); err != nil {
		s.Logger.Error("Failed to save customer to db: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, "Customer added successfully")
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

	// Get all customers from the db

	customers, err := s.db.GetCustomers(ctx)
	if err != nil {
		s.Logger.Error("Failed to get customers from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, customers)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, "Customer updated successfully")
}
