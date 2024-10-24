package testutils

import (
	"backend/model"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5432 user=dawn_user password=dawn_password dbname=dawn_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	err = db.AutoMigrate(&model.Post{})
	require.NoError(t, err)
	db.Exec("TRUNCATE TABLE users")
	db.Exec("TRUNCATE TABLE posts")

	return db
}
