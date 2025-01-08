package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
)

func (s *Server) AddCategory(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	category := models.Category{}

	// Unmarshal the request body into the category struct
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate category
	if err := category.Validate(); err != nil {
		s.Logger.Error("Failed to validate category: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// save to db
	if err := s.db.AddCategory(ctx, category); err != nil {
		s.Logger.Error("Failed to save category to db: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, "Category added successfully")
}

func (s *Server) GetCategories(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get all catgeories from the db

	categories, err := s.db.GetCategories(ctx)
	if err != nil {
		s.Logger.Error("Failed to get categories from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, categories)
}
