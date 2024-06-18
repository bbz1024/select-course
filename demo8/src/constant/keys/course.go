package keys

const (
	CourseHsetKey             = "course_hset:%d"
	CourseCapacityKey         = "course_capacity"
	CourseCategoryIDKey       = "course_category_id"
	CourseScheduleDurationKey = "course_schedule_duration"
	CourseScheduleWeekKey     = "course_schedule_week"
	CourseSequenceKey         = "course_sequence" // 课程时间序列号，确保执行有序
)
