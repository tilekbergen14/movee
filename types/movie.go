package types

type Movie struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewMovie(name, description, image string) *Movie {
	return &Movie{
		Name:        name,
		Description: description,
		Image:       image,
	}
}
