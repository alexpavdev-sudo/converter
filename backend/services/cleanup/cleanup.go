package cleanup

import (
	"context"
	"converter/app"
	"converter/config"
	"converter/entities"
	"log"
	"os"
	"time"
)

func Start(ctx context.Context) {
	interval := config.CleanupInterval * time.Second
	duration := config.SessionDuration * time.Second
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		cleanOldGuests(duration)

		for {
			select {
			case <-ticker.C:
				cleanOldGuests(duration)
			case <-ctx.Done():
				log.Println("Cleanup stopped:", ctx.Err())
				return
			}
		}
	}()

	log.Printf("Cleanup started: check every %v, delete older than %v", interval, duration)
}

func cleanOldGuests(duration time.Duration) {
	log.Println("Cleanup start")
	batchSize := 10
	var lastID uint = 0

	for {
		guests, err := nextGuests(duration, lastID, batchSize)
		if err != nil {
			log.Println("Ошибка:", err)
		}

		if len(guests) == 0 {
			break
		}

		for _, guest := range guests {
			cleanGuest(guest, duration)
		}

		lastID = guests[len(guests)-1].ID

		if len(guests) < batchSize {
			break
		}
	}
}

func nextGuests(duration time.Duration, lastID uint, batchSize int) ([]entities.Guest, error) {
	var guests []entities.Guest

	thresholdTime := time.Now().Add(-duration)

	err := app.App().DB.Where("created_at < ? AND id > ?", thresholdTime, lastID).
		Order("id ASC").
		Limit(batchSize).
		Find(&guests).Error

	return guests, err
}

func cleanGuest(guest entities.Guest, duration time.Duration) {
	thresholdTime := time.Now().Add(-duration)
	var files []entities.File
	err := app.App().DB.Model(&entities.File{}).
		Select("files.*").
		Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
		Where("guest_files.guest_id = ? AND files.created_at < ?", guest.ID, thresholdTime).
		Find(&files).Error
	if err != nil {
		log.Println("Ошибка db:", err)
		return
	}

	for _, file := range files {
		err = deleteFile(file)
		if err != nil {
			log.Println("Ошибка удаления файла:", err)
		}
	}
	app.ClearCache(guest.ID)
	deleteGuest(guest)
}

func deleteGuest(guest entities.Guest) {
	tx := app.App().StartTransaction()
	if tx.Error != nil {
		log.Println("Ошибка db:", tx.Error)
	}
	defer tx.Rollback()

	count, err := app.App().FileRepo.GetCountFiles(guest.ID)
	if err != nil {
		log.Println("Ошибка db:", err)
	}
	if count > 0 {
		return
	}
	if err := app.App().DB.Where("id = ?", guest.ID).Delete(&guest).Error; err != nil {
		log.Printf("Ошибка удаления гостя %d", guest.ID)
	}
	if err := os.Remove(guest.PersonalPath(true)); err != nil {
		log.Println("Ошибка удаления папки гостя:", err)
	}
	if err := tx.Commit().Error; err != nil {
		log.Println("Ошибка db:", err)
	}
	log.Printf("Удален гость %d", guest.ID)
}

func deleteFile(file entities.File) error {
	tx := app.App().StartTransaction()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Where("id = ?", file.ID).Delete(&file).Error; err != nil {
		return err
	}
	if err := os.Remove(file.PathFull()); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	log.Printf("Удален файл %d по истечению времени хранения", file.ID)

	return nil
}
