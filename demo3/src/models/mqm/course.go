package mqm

type ConsumerType uint
type CourseReq struct {
	UserID    uint         `json:"user_id"`
	CourseID  uint         `json:"course_id"`
	Type      ConsumerType `json:"type"`
	CreatedAt int64        `json:"created_at"`
}

const (
	SelectType ConsumerType = iota
	BackType
)
