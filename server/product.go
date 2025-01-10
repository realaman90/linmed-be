package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
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

func (s *Server) AddProduct(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	product := models.Product{}

	// Unmarshal the request body into the product struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate product
	if err := product.Validate(); err != nil {
		s.Logger.Error("Failed to validate product: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// save to db
	id, err := s.db.AddProduct(ctx, product)
	if err != nil {
		s.Logger.Error("Failed to save product to db: ", err)
		writeJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"id":      id,
		"message": "Product added successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the product id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("Product id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Product id is required")
		return
	}

	// Get the product from the db
	product, err := s.db.GetProduct(ctx, id)
	if err != nil {
		s.Logger.Error("Failed to get product from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the product
	writeJSONResponse(w, http.StatusOK, product)
}

func (s *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// Get all products from the db

	products, err := s.db.GetProducts(ctx)
	if err != nil {
		s.Logger.Error("Failed to get products from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, products)
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the product id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("Product id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Product id is required")
		return
	}

	product := models.Product{}

	// Unmarshal the request body into the product struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate product
	if err := product.Validate(); err != nil {
		s.Logger.Error("Failed to validate product: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// update product in db
	if err := s.db.UpdateProduct(ctx, product); err != nil {
		s.Logger.Error("Failed to update product in db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, "Product updated successfully")
}

func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the product id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("Product id is required")
		writeJSONResponse(w, http.StatusBadRequest, "Product id is required")
		return
	}

	// delete product from db
	if err := s.db.DeleteProduct(ctx, id); err != nil {
		s.Logger.Error("Failed to delete product from db: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	writeJSONResponse(w, http.StatusOK, "Product deleted successfully")
}
