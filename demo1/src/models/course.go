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
	// 例如： 星期一 上午 08:10 ~ 11:50 | 星期二 下午 14:10 ~ 16:50 | 星期三 晚上 18:50 ~ 21:20
	Duration string `json:"duration" gorm:"type:varchar(32);not null;comment:上课时间段"`
}
