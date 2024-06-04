package models

type BaseModel struct {
	ID uint `gorm:"primarykey" json:"id"`
}
