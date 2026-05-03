package entities

import "time"

type Notification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Detail    string     `gorm:"column:detail;type:text;not null" json:"detail"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	IsSend    uint8      `gorm:"column:is_send;default:0;not null" json:"is_send"`
	Type      TypeNotify `gorm:"column:type;not null" json:"type"`
	GuestID   uint       `gorm:"column:guest_id;type:bigint;not null" json:"guest_id"`
}

func (Notification) TableName() string {
	return "notifications"
}

type TypeNotify uint8

const (
	System TypeNotify = iota + 1
	User
)

func (t TypeNotify) String() string {
	switch t {
	case System:
		return "системное"
	case User:
		return "пользовательское"
	default:
		return "unknown"
	}
}
