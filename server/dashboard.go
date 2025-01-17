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

func (s *Server) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// get query month, week, day
	period := r.URL.Query().Get("period")

	tasks, err := s.db.GetUpcomingStaionWithTask(ctx, period)
	if err != nil {
		s.Logger.Error("Failed to get tasks from db: ", err)
		errorResposne(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tasks)
}
