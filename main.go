package main

import "log"

func main() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewTodoRepository(db) // to be passed to the handlers
}
