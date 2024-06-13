package mock

import (
	"context"
	"fmt"
	"select-course/demo4/src/constant/keys"
	"select-course/demo4/src/models"
	"select-course/demo4/src/storage/cache"
	"select-course/demo4/src/storage/database"
	"strconv"
	"testing"
)

func TestInsertUsers(t *testing.T) {
	var users []models.User
	for i := 1; i <= 100; i++ {
		users = append(users, models.User{
			BaseModel: models.BaseModel{ID: uint(i)},
			UserName:  "users" + strconv.Itoa(i),
			Password:  "password" + strconv.Itoa(i),
		})
	}
	// 清空
	if err := database.Client.Where("1=1").Delete(&models.User{}).Error; err != nil {
		t.Error(err)
	}
	if err := database.Client.Create(&users).Error; err != nil {
		t.Error(err)
	}
}
func TestInsertUsersRedis(t *testing.T) {
	var users []models.User
	if err := database.Client.Find(&users).Error; err != nil {
		t.Error(err)
	}
	for _, user := range users {
		err := cache.RDB.HSet(context.Background(),
			fmt.Sprintf(keys.UserInfoHsetKey, user.ID),
			keys.UserNameKey, user.UserName,
			keys.UserFlagKey, uint16(user.Flag),
		).Err()
		if err != nil {
			t.Error(err)
		}
	}
}
