package server

import (
	"context"
	"net/http"
)

func (s *Server) GetAllNumbers(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	numbers, err := s.db.GetAllNumbers(ctx)
	if err != nil {
		s.Logger.Error("Failed to get numbers from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, numbers)
}

func (s *Server) GetExpiryTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startDate := r.URL.Query().Get("startDate")
	if startDate == "" {
		errorResposne(w, http.StatusBadRequest, "startDate is required")
		return
	}

	endDate := r.URL.Query().Get("endDate")
	if endDate == "" {
		errorResposne(w, http.StatusBadRequest, "endDate is required")
		return
	}

	customerId := r.URL.Query().Get("customerId")

	// Get pagination parameters
	page, limit := s.validatePageLimit(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))

	tasks, total, err := s.db.GetExpiringProducts(ctx, startDate, endDate, customerId, page, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch expiring products")
		errorResposne(w, http.StatusInternalServerError, "Failed to fetch expiring products")
		return
	}

	res := paginatedResponse{
		Total: total,
		Data:  tasks,
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetInspectionTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startDate := r.URL.Query().Get("startDate")
	if startDate == "" {
		errorResposne(w, http.StatusBadRequest, "startDate is required")
		return
	}

	endDate := r.URL.Query().Get("endDate")
	if endDate == "" {
		errorResposne(w, http.StatusBadRequest, "endDate is required")
		return
	}

	// Get pagination parameters
	page, limit := s.validatePageLimit(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))

	tasks, total, err := s.db.GetInspectionTasks(ctx, startDate, endDate, page, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch inspection tasks")
		errorResposne(w, http.StatusInternalServerError, "Failed to fetch inspection tasks")
		return
	}

	res := paginatedResponse{
		Total: total,
		Data:  tasks,
	}

	writeJSONResponse(w, http.StatusOK, res)
}
