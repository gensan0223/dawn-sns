package handler

import (
	"backend/model"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func GetUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []model.User
		if result := db.Find(&users); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func CreateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestUser model.User
		err := json.NewDecoder(r.Body).Decode(&requestUser)
		if err != nil {
			http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
			return
		}

		if requestUser.Name == "" || requestUser.Email == "" || requestUser.Password == "" {
			http.Error(w, "不正な入力項目があります", http.StatusBadRequest)
			return
		}
		user := model.User{
			Name:     requestUser.Name,
			Email:    requestUser.Email,
			Password: requestUser.Password,
		}
		db.Create(&user)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "リクエスト作成に失敗しました", http.StatusInternalServerError)
			return
		}
	}
}
