package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aakash-tyagi/linmed/models"
	"github.com/gorilla/mux"
)

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	user := models.User{}

	// Unmarshal the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate user
	if err := user.Validate(); err != nil {
		s.Logger.Error("Failed to validate user: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// save to db
	id, err := s.db.AddUser(ctx, user)
	if err != nil {
		s.Logger.Error("Failed to save user to db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Id      uint   `json:"id"`
		Message string `json:"message"`
	}{
		Id:      id,
		Message: "User added successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the user id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("User id is required")
		writeJSONResponse(w, http.StatusBadRequest, "User id is required")
		return
	}

	// Get the user from the db
	user, err := s.db.GetUser(ctx, id)
	if err != nil {
		s.Logger.Error("Failed to get user from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return the user
	writeJSONResponse(w, http.StatusOK, user)
}

// update user

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Get the user id from the request
	id := mux.Vars(r)["id"]

	// validate id
	if id == "" {
		s.Logger.Error("User id is required")
		writeJSONResponse(w, http.StatusBadRequest, "User id is required")
		return
	}

	user := models.User{}

	// Unmarshal the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.Logger.Error("Failed to decode request body: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate user
	if err := user.Validate(); err != nil {
		s.Logger.Error("Failed to validate user: ", err)
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// update user in db
	if err := s.db.UpdateUser(ctx, id, user); err != nil {
		s.Logger.Error("Failed to update user in db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	ID, _ := strconv.Atoi(id)

	res := struct {
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		Id:      ID,
		Message: "User updated successfully",
	}

	// return success
	writeJSONResponse(w, http.StatusOK, res)

}

// get users

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, limitInt := s.validatePageLimit(page, limit)

	// Get the users from the db
	users, totalUsers, err := s.db.GetUsers(ctx, pageInt, limitInt)
	if err != nil {
		s.Logger.Error("Failed to get users from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := paginatedResponse{
		Total: totalUsers,
		Data:  users,
	}
	// return the users
	writeJSONResponse(w, http.StatusOK, res)
}
