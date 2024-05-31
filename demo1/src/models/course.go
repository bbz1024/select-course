package models

import "time"

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

	// 上课周 1 ~ 5
	Week uint8 `json:"week" gorm:"type:int;"`
	// 例如：  上午 08:10 ~ 11:50 |  下午 14:10 ~ 16:50 |  晚上 18:50 ~ 21:20
	Duration string `json:"duration" gorm:"type:varchar(32);not null;comment:上课时间段"`
	// 容纳人数
	Capacity uint `json:"capacity" gorm:"type:int;not null;comment:容纳人数"`
}

// UserCourse 用户选课关系表
type UserCourse struct {
	UserID    uint      `json:"userID" gorm:"not null;uniqueIndex:user_course;comment:用户ID"`
	User      *User     `json:"users" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:用户"`
	CourseID  uint      `json:"courseID" gorm:"not null;uniqueIndex:user_course;comment:课程ID"`
	Course    *Course   `json:"course" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:课程"`
	CreatedAt time.Time `json:"created_at" comment:"创建时间"`
}
