package mqm

type ConsumerType uint
type CourseReq struct {
	UserID   uint         `json:"user_id"`
	CourseID uint         `json:"course_id"`
	Type     ConsumerType `json:"type"`
}

const (
	SelectType ConsumerType = iota
	BackType
)
