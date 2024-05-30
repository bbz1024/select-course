package models

type User struct {
	BaseModel
	UserName string `json:"username" gorm:"type:varchar(64);not null;index;comment:用户名称;"`
	Password string `json:"password" gorm:"type:varchar(64);not null;comment:密码;"`
}
