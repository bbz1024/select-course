package models

type Flag uint32
type User struct {
	BaseModel
	UserName string `json:"username" gorm:"type:varchar(64);not null;index;comment:用户名称;"`
	Password string `json:"password" gorm:"type:varchar(64);not null;comment:密码;"`
	Flag     Flag   `json:"flag" gorm:"type:int;not null;comment:用户标准位记录着选课的已选字段;"`
}

// SetBit 设置指定位置的位为1，表示时间槽被占用
func (tb *Flag) SetBit(slot int) {
	if slot < 0 || slot >= 16 {
		panic("Slot out of range")
	}
	*tb |= 1 << slot
}

// ClearBit 清除指定位置的位，表示释放时间槽
func (tb *Flag) ClearBit(slot int) {
	if slot < 0 || slot >= 16 {
		panic("Slot out of range")
	}
	*tb &= ^(1 << slot)
}

// TestBit 检查指定位置的位是否为1，即时间槽是否被占用
func (tb *Flag) TestBit(slot int) bool {
	if slot < 0 || slot >= 16 {
		panic("Slot out of range")
	}
	return (*tb>>(slot))&1 == 1
}
