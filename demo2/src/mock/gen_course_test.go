package mock

import (
	"context"
	"fmt"
	"math/rand"
	"select-course/demo2/src/constant/keys"
	"select-course/demo2/src/models"
	"select-course/demo2/src/storage/cache"
	"select-course/demo2/src/storage/database"
	"testing"
)

func TestInsertCourse(t *testing.T) {
	// 星期一 上午 08:10 ~ 11:50 | 星期二 下午 14:10 ~ 16:50 | 星期三 晚上 18:50 ~ 21:20
	var course []models.Course
	for i := 1; i < 11; i++ {
		course = append(course, models.Course{
			BaseModel: models.BaseModel{
				ID: uint(i),
			},
			Title:      fmt.Sprintf("课程%d", i),
			CategoryID: uint(rand.Intn(5) + 1),
			ScheduleID: uint(rand.Intn(15) + 1),
			Capacity:   uint(rand.Intn(10) + 1),
		})
	}
	database.Client.Create(&course)

}
func TestInsertSchedule(t *testing.T) {
	var schedule []models.Schedule
	week := 5
	duration := 3
	for i := 0; i < week; i++ {
		for j := 1; j <= duration; j++ {
			schedule = append(schedule, models.Schedule{
				Week:     models.Week(i),
				Duration: models.Duration(j),
			})
		}
	}
	fmt.Println(schedule)
	database.Client.Create(&schedule)

}
func TestInsertCourseCategory(t *testing.T) {
	var courseCategory []models.CourseCategory
	for i := 1; i < 6; i++ {
		courseCategory = append(courseCategory, models.CourseCategory{
			Name: fmt.Sprintf("分类%d", i),
			BaseModel: models.BaseModel{
				ID: uint(i),
			},
		})
	}
	database.Client.Create(&courseCategory)
}

// 预热到redis
func TestPreheatMysql2Redis(t *testing.T) {
	var course []models.Course
	if err := database.Client.
		Preload("Schedule").
		Preload("Category").
		Find(&course).Error; err != nil {
		t.Error(err)
	}
	for _, v := range course {
		if err := cache.RDB.HSet(
			context.Background(),
			fmt.Sprintf(keys.CourseHsetKey, v.ID),
			keys.CourseCategoryIDKey, v.CategoryID,
			keys.CourseScheduleDurationKey, uint(v.Schedule.Duration),
			keys.CourseScheduleWeekKey, uint(v.Schedule.Week),
			keys.CourseCapacityKey, v.Capacity,
		).Err(); err != nil {
			t.Error(err)
		}
	}
}
