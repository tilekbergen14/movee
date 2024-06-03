package storage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/tilekbergen14/movee-back/types"
)

type Storage interface {
	CreateMovie(*types.Movie) error
	GetMovies() ([]*types.Movie, error)
}

type PostgresSql struct {
	db *sql.DB
}

func (s *PostgresSql) Init() error {
	return s.CreateMovieTable()
}

func (s *PostgresSql) CreateMovieTable() error {
	query := `create table if not exists movie(
		id serial primary key,
		name varchar(250),
		description text,
		image text
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func NewPostgresStorage() (*PostgresSql, error) {
	connStr := "user=postgres dbname=movee password=tqbank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresSql{
		db: db,
	}, nil
}

func (s *PostgresSql) CreateMovie(movie *types.Movie) error {
	query := `insert into movie (name, description, image)
	values ($1, $2, $3)
	`
	_, err := s.db.Exec(query, movie.Name, movie.Description, movie.Image)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresSql) GetMovies() ([]*types.Movie, error) {
	query := `select * from movie`
	res, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	movies := []*types.Movie{}

	for res.Next() {
		movie := &types.Movie{}
		if err := res.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.Image,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
