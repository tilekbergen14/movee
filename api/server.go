package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tilekbergen14/movee-back/types"
)

type Server struct {
	ListenAddr string
}

func NewServer(Addr string) *Server {
	return &Server{
		ListenAddr: Addr,
	}
}

func (s *Server) Run() error {

	router := mux.NewRouter()
	router.HandleFunc("/", makeHttpHandleFunc(s.handleHome))

	return http.ListenAndServe(s.ListenAddr, router)
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, 500, "Something went wrong!")
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) error {
	movies := [10]*types.Movie{}
	movies[0] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[1] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[2] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[3] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[4] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[5] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[6] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[7] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[8] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")
	movies[9] = types.NewMovie("Attack on Titan", "Some long description", "imagedata")

	return writeJSON(w, http.StatusOK, movies)
}
