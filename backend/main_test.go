package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5555 user=postgres password=password dbname=test_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&User{})
	require.NoError(t, err)

	return db
}

func CreateTestUsers(t *testing.T, db *gorm.DB) {
	bio := "bio example"
	bio2 := "bio example2"
	users := []User{
		{
			Name:     "Gem",
			Email:    "example@mail.com",
			Password: "test",
			Bio:      &bio,
		},
		{
			Name:     "Mma",
			Email:    "example2@mail.com",
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

	rr := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var users []User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Gen", users[0].Name)
	require.Equal(t, "Mma", users[1].Name)
}
