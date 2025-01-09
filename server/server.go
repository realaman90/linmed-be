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

	r.HandleFunc("/api/v1/user", s.AddUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", s.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/user/{id}", s.UpdateUser).Methods("PUT")

	r.HandleFunc("/api/v1/category", s.AddCategory).Methods("POST")
	r.HandleFunc("/api/v1/category/{id}", s.GetCategories).Methods("GET")

	// r.HandleFunc("/api/v1/product", s.AddProduct).Methods("POST")
	// r.HandleFunc("/api/v1/product/{id}", s.GetProduct).Methods("GET")
	// r.HandleFunc("/api/v1/product/{id}", s.UpdateProduct).Methods("PUT")
	// r.HandleFunc("/api/v1/product/{id}", s.DeleteProduct).Methods("DELETE")
	// r.HandleFunc("/api/v1/product", s.GetProducts).Methods("GET")

}

func (s *Server) Start() {
	r := mux.NewRouter()

	s.RegisterRoutes(r)
	s.corsMiddleware(r)

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
