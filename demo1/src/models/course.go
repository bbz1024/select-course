package models

// CourseCategory 课程种类
type CourseCategory struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar(32);not null;comment:分类名称"`
}

// Course 课程
type Course struct {
	BaseModel
	Title string `json:"title" gorm:"type:varchar(64);not null;comment:课程名称"`
	// 课程分类
	CategoryID uint            `json:"categoryID" gorm:"not null;comment:分类ID"`
	Category   *CourseCategory `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:课程分类"`
	//	上课时间段 早中晚
	ScheduleID uint      `json:"ScheduleID" gorm:"not null;comment:分类ID"`
	Schedule   *Schedule `json:"schedule" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:课程时间"`
}

type Duration uint8
type Week uint

const (
	Morning   Duration = iota // 上午
	AfterNoon                 // 下午
	Evening                   // 晚上

)
const (
	// 上课周
	_ Week = iota
	Mon
	Tue
	Wed
	Thu
	Fri
)

// Schedule 时刻表
type Schedule struct {
	BaseModel
	Duration Duration `json:"duration" gorm:"type:int;"`
	//  上课周 周一到周五
	Week Week `json:"week" gorm:"type:int;"`
}

/*
// UserCourse 用户选课关系表
type UserCourse struct {
	UserID    uint      `json:"userID" gorm:"not null;uniqueIndex:user_course;comment:用户ID"`
	User      *User     `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:用户"`
	CourseID  uint      `json:"courseID" gorm:"not null;uniqueIndex:user_course;comment:课程ID"`
	Course    *Course   `json:"course" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:课程"`
	CreatedAt time.Time `json:"created_at" comment:"创建时间"`
}

*/
