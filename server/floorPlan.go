package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
)

func (s *Server) AddFloorPlan(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	floorPlan := models.FloorPlan{}

	// get cusoemr id from the request
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// custoemr id to uint
	customerIDUint, err := s.stringToUint(customerID)
	if err != nil {
		s.Logger.Error("Failed to convert customer id to int: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the request body into the floorPlan struct
	if err := json.NewDecoder(r.Body).Decode(&floorPlan); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// set customer id
	floorPlan.CustomerID = uint(customerIDUint)

	// validate floorPlan
	if err := floorPlan.Validate(); err != nil {
		s.Logger.Error("Failed to validate floorPlan: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// save to db
	id, err := s.db.AddFloorPlan(ctx, floorPlan)
	if err != nil {
		s.Logger.Error("Failed to save floorPlan to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      id,
		"message": "FloorPlan added successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)

}

func (s *Server) GetFloorPlans(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get the customer id from the request
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// custoemr id to uint
	customerIDUint, err := s.stringToUint(customerID)
	if err != nil {
		s.Logger.Error("Failed to convert customer id to int: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the page and limit query parameters

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// validate page and limit
	pageInt, limitInt := s.validatePageLimit(page, limit)

	// Get the floorPlans from the db
	floorPlans, total, err := s.db.GetFloorPlans(ctx, customerIDUint, pageInt, limitInt)
	if err != nil {
		s.Logger.Error("Failed to get floorPlans from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	// paginated response
	res := paginatedResponse{
		Total: total,
		Data:  floorPlans,
	}

	// return the floorPlans
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetFloorPlan(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get the floorPlan id from the request
	id := mux.Vars(r)["floorPlanID"]

	// validate id
	if id == "" {
		s.Logger.Error("FloorPlan id is required")
		writeJSONResponse(w, http.StatusBadRequest, "FloorPlan id is required")
		return
	}

	// Get the floorPlan from the db
	floorPlan, err := s.db.GetFloorPlan(ctx, id)
	if err != nil {
		s.Logger.Error("Failed to get floorPlan from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return the floorPlan
	writeJSONResponse(w, http.StatusOK, floorPlan)
}

func (s *Server) UpdateFloorPlan(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get the floorPlan id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("FloorPlan id is required")
		writeJSONResponse(w, http.StatusBadRequest, "FloorPlan id is required")
		return
	}

	floorPlan := models.FloorPlan{}

	// Unmarshal the request body into the floorPlan struct
	if err := json.NewDecoder(r.Body).Decode(&floorPlan); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate floorPlan
	if err := floorPlan.Validate(); err != nil {
		s.Logger.Error("Failed to validate floorPlan: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// update floorPlan in db
	if err := s.db.UpdateFloorPlan(ctx, id, floorPlan); err != nil {
		s.Logger.Error("Failed to update floorPlan in db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	ID, _ := s.stringToUint(id)

	res := struct {
		Id      uint   `json:"id"`
		Message string `json:"message"`
	}{
		Id:      ID,
		Message: "FloorPlan updated successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) DeleteFloorPlan(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get the floorPlan id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("FloorPlan id is required")
		writeJSONResponse(w, http.StatusBadRequest, "FloorPlan id is required")
		return
	}

	// delete floorPlan from db
	if err := s.db.DeleteFloorPlan(ctx, id); err != nil {
		s.Logger.Error("Failed to delete floorPlan from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]string{
		"id":      id,
		"message": "FloorPlan deleted successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}
