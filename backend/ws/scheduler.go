package ws

import (
	ws_dto "converter/dto/ws"
	"converter/entities"
	"log"
	"time"

	"gorm.io/gorm"
)

func StartNotificationsWatcher(db *gorm.DB, hub *Hub) {
	for {
		log.Printf("count websocket guests: %d", hub.countGuest())
		log.Printf("count websocket connections: %d", hub.totalConnections())

		var notifications []entities.Notification
		err := db.Where("is_send = 0").
			Order("id asc").
			Limit(100).
			Find(&notifications).Error
		if err != nil {
			log.Printf("notifications watcher error db: %v", err)
		}

		notificationsDto := ws_dto.NotificationToDTO().MapSlice(notifications)
		ids := make([]uint, 0, len(notificationsDto))

		for _, notificationDto := range notificationsDto {
			ids = append(ids, notificationDto.ID)
			select {
			case hub.Notify <- notificationDto:
			default:
				log.Printf("notify channel full")
				time.Sleep(10 * time.Second)
			}
		}

		if len(ids) > 0 {
			err = db.Model(&entities.Notification{}).
				Where("id IN ?", ids).
				Update("is_send", 1).Error
			if err != nil {
				log.Printf("error update is_send error for notifications: %v", err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
