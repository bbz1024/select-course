package mock

import (
	"select-course/demo1/src/models"
	"select-course/demo1/src/storage/database"
	"strconv"
	"testing"
)

func TestInsertUsers(t *testing.T) {
	var users []models.User
	for i := 0; i < 100; i++ {
		users = append(users, models.User{
			BaseModel: models.BaseModel{ID: uint(i)},
			UserName:  "user" + strconv.Itoa(i),
			Password:  "password" + strconv.Itoa(i),
		})
	}
	database.Client.Create(users)
}
