package main

import (
	"log"

	"github.com/tilekbergen14/movee-back/api"
)

func main() {
	s := api.NewServer(":4000")
	log.Println("Server running on port ", s.ListenAddr)
	log.Println(s.Run())
}
