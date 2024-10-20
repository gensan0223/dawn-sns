package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5432 user=dawn_user password=dawn_password dbname=dawn_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&User{})
	err = db.AutoMigrate(&Post{})
	require.NoError(t, err)

	return db
}

func CreateTestUsers(t *testing.T, db *gorm.DB) {
	bio := "bio example"
	bio2 := "bio example2"
	users := []User{
		{
			Name:     "Gem",
			Email:    "example1@mail.com",
			Password: "test",
			Bio:      &bio,
		},
		{
			Name:     "Mma",
			Email:    "example8@mail.com",
			Password: "test2",
			Bio:      &bio2,
		},
	}
	for _, user := range users {
		err := db.Create(&user).Error
		require.NoError(t, err)
	}
}

func TestGetUsersHandler(t *testing.T) {
	db := setupTestDB(t)
	CreateTestUsers(t, db)

	req, err := http.NewRequest("GET", "/users", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var users []User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Gem", users[0].Name)
	require.Equal(t, "Mma", users[1].Name)
}
func TestCreatePostHandler(t *testing.T) {
	testPost := Post{
		UserID:  1,
		Content: "test content",
	}
	reqBody, _ := json.Marshal(testPost)
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	db := setupTestDB(t)
	handler := CreatePostHandler(db)
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var createdPost Post
	err = json.NewDecoder(rec.Body).Decode(&createdPost)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	require.Equal(t, testPost.UserID, createdPost.UserID)
	require.Equal(t, testPost.Content, createdPost.Content)
}
