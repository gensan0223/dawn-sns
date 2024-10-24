package handler

import (
	"backend/model"
	"backend/testutils"
	"backend/util"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func CreateTestUser(t *testing.T, db *gorm.DB) (*model.User, error) {
	user := model.User{
		Name:     util.RandomString(5),
		Email:    util.RandomEmail(5),
		Password: util.RandomString(10),
		Bio:      util.RandomBio(20),
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	require.NoError(t, err)
	return &user, nil
}

func CreateTestUsers(t *testing.T, count int, db *gorm.DB) []model.User {
	var users []model.User
	for i := 0; i < count; i++ {
		user, _ := CreateTestUser(t, db)
		users = append(users, *user)
	}
	return users
}

func TestGetUserHandler(t *testing.T) {
	db := testutils.SetupTestDB(t)

	testUser, _ := CreateTestUser(t, db)
	idStr := strconv.FormatUint(uint64(testUser.ID), 10)
	getUserHandlerURL := "/users/" + idStr

	req, err := http.NewRequest("GET", getUserHandlerURL, nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	handler := GetUserHandler(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var user model.User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	require.Equal(t, testUser.Name, user.Name)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.CreatedAt, user.CreatedAt)
}

func TestGetUsersHandler(t *testing.T) {
	db := testutils.SetupTestDB(t)

	testUsers := CreateTestUsers(t, 2, db)

	req, err := http.NewRequest("GET", "/users", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var users []model.User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	require.NoError(t, err)
	require.Equal(t, testUsers[0].Name, users[0].Name)
	require.Equal(t, testUsers[1].Email, users[1].Email)
	require.Equal(t, testUsers[1].CreatedAt, users[1].CreatedAt)
}

func TestCreateUserHandler(t *testing.T) {
	user := model.User{
		Name:     util.RandomString(5),
		Email:    util.RandomEmail(8),
		Password: util.RandomString(10),
		Bio:      util.RandomBio(20),
	}
	reqBody, err := json.Marshal(&user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	db := testutils.SetupTestDB(t)
	handler := CreateUserHandler(db)
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var createdUser model.User
	err = json.NewDecoder(rec.Body).Decode(&createdUser)
	require.NoError(t, err)
	require.Equal(t, user.Name, createdUser.Name)
}
