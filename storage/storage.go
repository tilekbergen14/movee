package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/tilekbergen14/movee-back/types"
)

type Storage interface {
	CreateMovie(*types.Movie) error
	GetMovies() ([]*types.Movie, error)
	DeleteMovie(id int) error
	UpdateMovie(*types.Movie) error
	GetMovieByid(id int) (*types.Movie, error)
}

func (s *PostgresSql) DeleteMovie(id int) error {
	query := `DELETE FROM movie
				WHERE id = $1`
	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresSql) UpdateMovie(movie *types.Movie) error {
	query := `UPDATE movie
				SET name = $1, description = $2, image = $3
				WHERE id = $4;`
	_, err := s.db.Query(query, movie.Name, movie.Description, movie.Image, movie.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresSql) GetMovieByid(id int) (*types.Movie, error) {
	query := "select * from movie where id = $1"
	res, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	movie := &types.Movie{}
	empty := true
	for res.Next() {
		empty = false
		if err := res.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.Image,
		); err != nil {
			return nil, err
		}
	}
	if empty {
		return nil, fmt.Errorf("can't find movie")
	}
	return movie, nil
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

type PostgresSql struct {
	db *sql.DB
}
