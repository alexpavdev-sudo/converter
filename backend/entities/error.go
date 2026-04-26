package entities

import "time"

type Error struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID    uint      `gorm:"column:file_id;type:bigint;not null;index" json:"file_id"`
	Details   string    `gorm:"column:details;type:text" json:"details"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
}
