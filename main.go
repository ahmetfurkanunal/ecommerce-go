package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"ecommerce/handlers"
	"ecommerce/repository"

	_ "github.com/glebarez/sqlite"
)

func main() {
	dbPath := "ecommerce.db"

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("read schema:", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal("apply schema:", err)
	}

	sqlRepos := repository.NewSQLRepos(db)

	api := handlers.NewAPI(sqlRepos.Users, sqlRepos.Products, sqlRepos.Carts)

	http.HandleFunc("/users/register", api.HandleRegisterUser)
	http.HandleFunc("/users/login", api.HandleLogin)
	http.HandleFunc("/users/", api.HandleUpdateUser)
	http.HandleFunc("/users", api.HandleListUsers)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
