package ws_dto

import (
	"converter/components/mapper"
	"converter/entities"
	"time"
)

type MessageDto struct {
	Type    Type `json:"type"`
	Payload any  `json:"payload"`
}

type Type uint8

const (
	Notification Type = iota + 1
	RegisterAck
	Error
)

func (t Type) String() string {
	switch t {
	case Notification:
		return "Уведомление"
	case RegisterAck:
		return "Регистрация успешна"
	case Error:
		return "Ошибка"
	default:
		return "unknown"
	}
}

type NotificationDto struct {
	ID        uint                `json:"id"`
	Detail    string              `json:"detail"`
	GuestID   uint                `json:"guest_id"`
	Type      entities.TypeNotify `json:"type"`
	CreatedAt string              `json:"created_at"`
}

type RegisterAckDto struct {
	GuestID uint   `json:"guest_id"`
	Status  string `json:"status"`
}

type ErrorPayloadDto struct {
	Message string `json:"message"`
}

var notificationToDTO = mapper.New(func(n entities.Notification) NotificationDto {
	return NotificationDto{
		ID:        n.ID,
		Detail:    n.Detail,
		GuestID:   n.GuestID,
		Type:      n.Type,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
})

func NotificationToDTO() *mapper.Mapper[entities.Notification, NotificationDto] {
	return notificationToDTO
}
