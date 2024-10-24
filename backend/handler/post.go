package handler

import (
	"backend/model"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func CreatePostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPost model.Post
		err := json.NewDecoder(r.Body).Decode(&requestPost)
		if err != nil {
			http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
			return
		}

		if requestPost.Content == "" {
			http.Error(w, "コンテンツがありません", http.StatusBadRequest)
			return
		}

		post := model.Post{
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
