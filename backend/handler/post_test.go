package handler

import (
	"backend/model"
	"backend/testutils"
	"backend/util"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePostHandler(t *testing.T) {
	testPost := model.Post{
		UserID:  1,
		Content: util.RandomString(30),
	}
	reqBody, _ := json.Marshal(testPost)
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	db := testutils.SetupTestDB(t)
	handler := CreatePostHandler(db)
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var createdPost model.Post
	err = json.NewDecoder(rec.Body).Decode(&createdPost)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	require.Equal(t, testPost.UserID, createdPost.UserID)
	require.Equal(t, testPost.Content, createdPost.Content)
}
