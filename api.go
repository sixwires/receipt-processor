package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var Receipts = make(map[uuid.UUID]Receipt)

type APIServer struct {
	listenAddr string
}

// API constructor
func NewAPIServer(listenAddr string) *APIServer {
	// Return pointer to server w/ listen addr
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// Server start
func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Route handlers
	router.HandleFunc("/", makeHTTPHandleFunc(ping))
	router.PathPrefix("/receipts").Handler(s.ReceiptsHandler())

	// Open port and run server
	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// Health handler
func ping(w http.ResponseWriter, r *http.Request) error {
	return WriteJson(w, http.StatusOK, "success")
}

// Encode json values for responses
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Function signature for handler funcs
type ApiFunc func(http.ResponseWriter, *http.Request) error

// Api errors type
type ApiError struct {
	Error string
}

// Decorate handler functions into HTTPHandlerFunc type
func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	// Handle functions from our handlers
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle the error
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
