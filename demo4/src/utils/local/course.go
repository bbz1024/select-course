package local

import (
	"fmt"
	"select-course/demo4/src/models"
	"select-course/demo4/src/storage/database"
)

var CourseSchedule map[uint]*CourseScheduleModel

type CourseScheduleModel struct {
	Week     int
	Duration int
}

func init() {
	CourseSchedule = make(map[uint]*CourseScheduleModel)
	var courseList []*models.Course
	if err := database.Client.Model(&models.Course{}).Preload("Schedule").Find(&courseList).Error; err != nil {
		panic(err)
	}
	for _, course := range courseList {
		CourseSchedule[course.ID] = &CourseScheduleModel{
			Week:     int(course.Schedule.Week),
			Duration: int(course.Schedule.Duration),
		}
	}
}
func CalOffset(courseID uint) (offset int, err error) {
	model, ok := CourseSchedule[courseID]
	if !ok {
		return 0, fmt.Errorf("courseID %d not found", courseID)
	}
	offset = int(model.Week)*3 + model.Duration
	return offset, nil
}
