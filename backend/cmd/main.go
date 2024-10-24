package main

import (
	"backend/handler"
	"backend/model"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Follow{})
	println("test!!")
	return db
}

func setupRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", handler.GetUsersHandler(db))
	mux.HandleFunc("POST /users", handler.CreateUserHandler(db))
	mux.HandleFunc("POST /posts", handler.CreatePostHandler(db))
	return mux
}

func main() {
	db := initDB()
	mux := setupRoutes(db)
	log.Println("Server running on port 8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
