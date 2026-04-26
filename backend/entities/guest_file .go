package entities

type GuestFile struct {
	GuestID uint `gorm:"column:guest_id;primaryKey;autoIncrement:false;index" json:"guest_id"`
	FileID  uint `gorm:"column:file_id;primaryKey;autoIncrement:false;index" json:"file_id"`
}
