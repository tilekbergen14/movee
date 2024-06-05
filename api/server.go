package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

	fs := http.FileServer(http.Dir("./assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	router.Use(enableCORS)
	router.HandleFunc("/handlefile", makeHttpHandleFunc(s.handleFile))

	router.HandleFunc("/", makeHttpHandleFunc(s.handleHome))
	router.HandleFunc("/movies", makeHttpHandleFunc(s.handleMovies))
	router.HandleFunc("/movies/{id}", makeHttpHandleFunc(s.handleMoviesById))
	return http.ListenAndServe(s.ListenAddr, router)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "Yeah nice to meet you ahhahaah")
}

func (s *Server) handleMovies(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		res, err := s.Storage.GetMovies()
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, res)
	}
	if r.Method == "POST" {
		movie := &types.Movie{}

		if err := json.NewDecoder(r.Body).Decode(movie); err != nil {
			return err
		}
		if err := s.Storage.CreateMovie(movie); err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, "New movie created!")
	}
	return writeJSON(w, http.StatusMethodNotAllowed, "Method not allowed!")
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := os.OpenFile("./assets/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			return err
		}
		io.Copy(f, file)
		return writeJSON(w, http.StatusOK, "./assets/"+handler.Filename)
	}
	return writeJSON(w, http.StatusMethodNotAllowed, "Method not allowed!")

}

func (s *Server) handleMoviesById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("wrong id")
	}
	if r.Method == "GET" {
		movie, err := s.Storage.GetMovieByid(intId)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, movie)
	}
	if r.Method == "PATCH" {
		movie := &types.Movie{}
		if err := json.NewDecoder(r.Body).Decode(movie); err != nil {
			return err
		}
		movie.Id = intId
		err := s.Storage.UpdateMovie(movie)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, "Updated!")
	}
	if r.Method == "DELETE" {
		return s.Storage.DeleteMovie(intId)
	}
	return writeJSON(w, http.StatusMethodNotAllowed, "Method not allowed!")
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
