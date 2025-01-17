package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
)

func (s *Server) AddStation(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	station := models.Station{}
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	//get cusomter id from the request
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// get floor id and station id from the request
	floorId := mux.Vars(r)["floorPlanID"]
	if floorId == "" {
		s.Logger.Error("Floor id is required")
		errorResposne(w, http.StatusBadRequest, "Floor id is required")
		return
	}

	var floorIdUint uint
	// floor id to uint
	floorIdUint, err := s.stringToUint(floorId)
	if err != nil {
		s.Logger.Error("Failed to convert floor id to int: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// set customer id
	customerIDUint, err := s.stringToUint(customerID)
	if err != nil {
		s.Logger.Error("Failed to convert customer id to int: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	fId := uint(floorIdUint)

	station.CustomerID = uint(customerIDUint)
	station.FloorPlanID = &fId

	if err := station.Validate(); err != nil {
		s.Logger.Error("Failed to validate station: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := s.db.AddStation(ctx, station)
	if err != nil {
		s.Logger.Error("Failed to save station to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      id,
		"message": "Station added successfully",
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetStationById(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id := mux.Vars(r)["id"]
	if id == "" {
		s.Logger.Error("Station id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Station id is required")
		return
	}

	floorPlanId := mux.Vars(r)["floorPlanID"]
	if floorPlanId == "" {
		s.Logger.Error("Floor id is required")
		errorResposne(w, http.StatusBadRequest, "Floor id is required")
		return
	}

	// get station id
	stationId := mux.Vars(r)["stationID"]
	if stationId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	station, err := s.db.GetStation(ctx, stationId)
	if err != nil {
		s.Logger.Error("Failed to get station from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, station)
}

func (s *Server) UpdateStation(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	station := models.Station{}
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	//get cusomter id from the request
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// get floor id and station id from the request
	floorId := mux.Vars(r)["floorPlanID"]
	if floorId == "" {
		s.Logger.Error("Floor id is required")
		errorResposne(w, http.StatusBadRequest, "Floor id is required")
		return
	}

	stationId := mux.Vars(r)["stationID"]
	if stationId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	if err := station.Validate(); err != nil {
		s.Logger.Error("Failed to validate station: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	err := s.db.UpdateStation(ctx, stationId, station)
	if err != nil {
		s.Logger.Error("Failed to save station to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      stationId,
		"message": "Station updated successfully",
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) DeleteStation(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	stationId := mux.Vars(r)["stationID"]
	if stationId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	err := s.db.DeleteStation(ctx, stationId)
	if err != nil {
		s.Logger.Error("Failed to delete station from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      stationId,
		"message": "Station deleted successfully",
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetStations(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// validate page and limit query parameters
	pageInt, limitInt := s.validatePageLimit(page, limit)

	// Get the customer id from the request
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	// Get the floor id from the request
	floorId := mux.Vars(r)["floorPlanID"]
	if floorId == "" {
		s.Logger.Error("Floor id is required")
		errorResposne(w, http.StatusBadRequest, "Floor id is required")
		return
	}

	stations, total, err := s.db.GetStations(ctx, pageInt, limitInt, floorId, customerID)
	if err != nil {
		s.Logger.Error("Failed to get stations from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"total":    total,
		"stations": stations,
	}

	writeJSONResponse(w, http.StatusOK, res)
}
