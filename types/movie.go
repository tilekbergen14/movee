package types

import "math/rand"

type Movie struct {
	Id          int
	Name        string
	Description string
	Image       string
}

func NewMovie(name, description, image string) *Movie {
	return &Movie{
		Id:          rand.Intn(1000000000),
		Name:        name,
		Description: description,
		Image:       image,
	}
}
