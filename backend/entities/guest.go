package entities

import (
	"converter/config"
	"path/filepath"
	"time"
)

type Guest struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PersonalDir string    `gorm:"column:personal_dir;type:varchar(128);not null" json:"personal_dir"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
}

func (g *Guest) PersonalPath(isAbsolute bool) string {
	if g.PersonalDir == "" {
		return ""
	}
	if isAbsolute {
		return filepath.Join(config.UploadDir, config.PrefixGuestDir+g.PersonalDir)
	}
	return config.PrefixGuestDir + g.PersonalDir
}
