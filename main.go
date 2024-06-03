package main

import (
	"log"

	"github.com/tilekbergen14/movee-back/api"
	"github.com/tilekbergen14/movee-back/storage"
)

func main() {
	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	s := api.NewServer(":4000", store)
	log.Println("Server running on port ", s.ListenAddr)
	log.Println(s.Run())
}
