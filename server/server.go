package server

import (
	"net/http"

	"github.com/aakash-tyagi/linmed/config"
	database "github.com/aakash-tyagi/linmed/db"
	"github.com/gorilla/handlers"
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
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.HandleFunc("/", s.HealthCheck).Methods("GET")
	r.HandleFunc("/health", s.HealthCheck).Methods("GET")

	r.HandleFunc("/api/v1/user", s.AddUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", s.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/user/{id}", s.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/v1/users", s.GetUsers).Methods("GET")

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

	// s.corsMiddleware(r)
	s.RegisterRoutes(r)

	// Apply CORS middleware
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	s.Logger.Info("Starting server on port: ", s.Config.ServerPort)
	if err := http.ListenAndServe(":"+s.Config.ServerPort, handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(r)); err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info("Server stopped")
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {

	writeJSONResponse(w, 200, nil)
}
