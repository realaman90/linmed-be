package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakash-tyagi/linmed/models"
)

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {

	// ctx := context.TODO()

	user := models.User{}

	// Unmarshal the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(user)

}
