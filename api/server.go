package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tilekbergen14/movee-back/storage"
	"github.com/tilekbergen14/movee-back/types"
)

type Server struct {
	ListenAddr string
	Storage    storage.Storage
}

func NewServer(Addr string, storage storage.Storage) *Server {
	return &Server{
		ListenAddr: Addr,
		Storage:    storage,
	}
}

func (s *Server) Run() error {

	router := mux.NewRouter()
	router.HandleFunc("/", makeHttpHandleFunc(s.handleHome))
	router.HandleFunc("/movies", makeHttpHandleFunc(s.handleMovies))

	return http.ListenAndServe(s.ListenAddr, router)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "Yeah nice to meet you ahhahaah")
}

func (s *Server) handleMovies(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getAllMovies(w, r)
	}
	if r.Method == "POST" {
		movie := &types.Movie{}
		if err := json.NewDecoder(r.Body).Decode(movie); err != nil {
			return err
		}
		err := s.Storage.CreateMovie(movie)
		return err
	}
	return writeJSON(w, http.StatusOK, "New movie created!")
}

func (s *Server) getAllMovies(w http.ResponseWriter, r *http.Request) error {
	res, err := s.Storage.GetMovies()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, res)
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, 500, err.Error())
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
