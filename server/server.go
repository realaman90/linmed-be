package server

import (
	"net/http"

	"github.com/aakash-tyagi/linmed/config"
	database "github.com/aakash-tyagi/linmed/db"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Config *config.Config
	Logger *log.Logger
	db     *database.Database
}

func New(
	config *config.Config,
	log *log.Logger,
	db *database.Database,
) *Server {
	return &Server{
		Config: config,
		Logger: log,
		db:     db,
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", s.HealthCheck).Methods("GET")
	r.HandleFunc("/api/v1", s.AddUser).Methods("GET")
}

func (s *Server) Start() {
	r := mux.NewRouter()

	s.RegisterRoutes(r)

	http.Handle("/", r)
	s.Logger.Info("Starting server on port: ", s.Config.ServerPort)
	if err := http.ListenAndServe(":"+s.Config.ServerPort, nil); err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info("Server stopped")
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, 200, nil)
}
