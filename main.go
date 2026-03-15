package main

import (
	"log"
	"net/http"

	"practice5/db"
	"practice5/handler"
	"practice5/repository"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	if err := db.Seed(database); err != nil {
		log.Fatalf("failed to seed: %v", err)
	}

	repo := repository.NewRepository(database)
	h := handler.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.GetPaginatedUsers)
	mux.HandleFunc("/users/common-friends", h.GetCommonFriends)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
