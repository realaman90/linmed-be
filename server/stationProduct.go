package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
)

func (s *Server) AddStationProduct(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	device := models.StationProduct{}

	// Unmarshal the request body into the device struct
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate device
	if err := device.Validate(); err != nil {
		s.Logger.Error("Failed to validate device: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	// save to db
	id, err := s.db.AddStationProduct(ctx, device)
	if err != nil {
		s.Logger.Error("Failed to save device to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      id,
		"message": "Device added successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetStationProductById(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	stationProductId := mux.Vars(r)["id"]
	if stationProductId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	stationProduct, err := s.db.GetStationProductById(ctx, stationProductId)
	if err != nil {
		s.Logger.Error("Failed to get devices from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, stationProduct)
}

func (s *Server) UpdateStationProduct(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	stationProductId := mux.Vars(r)["id"]
	if stationProductId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	stationProduct := models.StationProduct{}
	if err := json.NewDecoder(r.Body).Decode(&stationProduct); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := stationProduct.Validate(); err != nil {
		s.Logger.Error("Failed to validate device: ", err)
		errorResposne(w, http.StatusBadRequest, err.Error())
		return
	}

	err := s.db.UpdateStationProduct(ctx, stationProductId, stationProduct)
	if err != nil {
		s.Logger.Error("Failed to update device: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      stationProductId,
		"message": "Device updated successfully",
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) DeleteStationProduct(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	stationProductId := mux.Vars(r)["id"]
	if stationProductId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	err := s.db.DeleteStationProduct(ctx, stationProductId)
	if err != nil {
		s.Logger.Error("Failed to delete device: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"id":      stationProductId,
		"message": "Device deleted successfully",
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetStationProducts(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")

	// validate page and limit query parameters
	pageInt, limitInt := s.validatePageLimit(page, limit)

	//get the station id and customer id from the request
	cusmtomerId := r.URL.Query().Get("customer_id")
	stationId := r.URL.Query().Get("station_id")

	if cusmtomerId == "" {
		s.Logger.Error("Customer id is required")
		errorResposne(w, http.StatusBadRequest, "Customer id is required")
		return
	}

	if stationId == "" {
		s.Logger.Error("Station id is required")
		errorResposne(w, http.StatusBadRequest, "Station id is required")
		return
	}

	stationProducts, total, err := s.db.GetStationProducts(ctx, pageInt, limitInt, cusmtomerId, stationId)
	if err != nil {
		s.Logger.Error("Failed to get devices from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"total": total,
		"data":  stationProducts,
	}

	writeJSONResponse(w, http.StatusOK, res)
}
