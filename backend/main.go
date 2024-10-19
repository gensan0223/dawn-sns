package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	Bio       *string   `gorm:"size:255"`
	Image     *string   `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	Content   string    `gorm:"size:400;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Follow struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null"`
	FollowerID uint      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func initDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Follow{})
	println("test!!")
	return db
}

func setupRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", GetUsersHandler(db))
	mux.HandleFunc("POST /posts", CreatePostHandler(db))
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

func GetUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []User
		if result := db.Find(&users); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func CreatePostHandler(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var requestPost Post
		err := json.NewDecoder(r.Body).Decode(&requestPost)

		if err != nil {
			http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
			return
		}

		if requestPost.Content == "" {
			http.Error(w, "コンテンツがありません", http.StatusBadRequest)
			return
		}

		post := Post{
			UserID:  requestPost.UserID,
			Content: requestPost.Content,
		}
		db.Create(&post)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			http.Error(w, "レスポンス作成に失敗しました", http.StatusInternalServerError)
			return
		}
	}
}
