package main

import (
	"backend/util"
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
	db.Exec("TRUNCATE TABLE users")
	db.Exec("TRUNCATE TABLE posts")

	return db
}

func CreateTestUser(t *testing.T, db *gorm.DB) (*User, error) {
	user := User{
		Name:     util.RandomString(5),
		Email:    util.RandomEmail(5),
		Password: util.RandomString(10),
		Bio:      util.RondomBio(20),
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	require.NoError(t, err)
	return &user, nil
}

func CreateTestUsers(t *testing.T, count int, db *gorm.DB) []User {
	var users []User
	for i := 0; i < count; i++ {
		user, _ := CreateTestUser(t, db)
		users = append(users, *user)
	}
	return users
}

func TestGetUsersHandler(t *testing.T) {
	db := setupTestDB(t)

	testUsers := CreateTestUsers(t, 2, db)

	req, err := http.NewRequest("GET", "/users", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var users []User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	require.NoError(t, err)
	require.Equal(t, testUsers[0].Name, users[0].Name)
	require.Equal(t, testUsers[1].Email, users[1].Email)
	require.Equal(t, testUsers[1].CreatedAt, users[1].CreatedAt)
}

func TestCreateUserHandler(t *testing.T) {
	user := User{
		Name:     util.RandomString(5),
		Email:    util.RandomEmail(8),
		Password: util.RandomString(10),
		Bio:      util.RondomBio(20),
	}
	reqBody, err := json.Marshal(&user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	db := setupTestDB(t)
	handler := CreateUserHandler(db)
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var createdUser User
	err = json.NewDecoder(rec.Body).Decode(&createdUser)
	require.NoError(t, err)
	require.Equal(t, user.Name, createdUser.Name)
}

func TestCreatePostHandler(t *testing.T) {
	testPost := Post{
		UserID:  1,
		Content: util.RandomString(30),
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
