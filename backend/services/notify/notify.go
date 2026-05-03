package notify

import (
	"converter/app"
	"converter/dto/web"
	"converter/entities"
	"encoding/json"
	"errors"
	"fmt"
)

type NotifyService struct {
}

func (s NotifyService) Notify(detail any, typeN entities.TypeNotify, guestID uint) error {
	detailAsBytes, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	n := &entities.Notification{
		Detail:  string(detailAsBytes),
		Type:    typeN,
		GuestID: guestID,
	}
	if err := app.App().DB.Create(n).Error; err != nil {
		return err
	}
	return nil
}

func (s NotifyService) NotifyConvertFileSuccess(fileId uint, originalName string) error {
	var guestFile entities.GuestFile
	err := app.App().DB.Model(&entities.GuestFile{}).
		Where("file_id = ?", fileId).
		First(&guestFile).Error
	if err != nil {
		return errors.New("error fetch guest_id")
	}

	detail := web.ResponseDto{
		Success: true,
		Data:    fmt.Sprintf("Файла %s обработан успешно", originalName),
	}

	err = s.Notify(
		detail,
		entities.User,
		guestFile.GuestID,
	)
	if err != nil {
		return fmt.Errorf("error save notification: %s", err.Error())
	}

	return nil
}

func (s NotifyService) NotifyConvertFileError(fileId uint, originalName string, error string) error {
	var guestFile entities.GuestFile
	err := app.App().DB.Model(&entities.GuestFile{}).
		Where("file_id = ?", fileId).
		First(&guestFile).Error
	if err != nil {
		return errors.New("error fetch guest_id")
	}

	detail := web.ResponseDto{
		Success: false,
		Data:    fmt.Sprintf("Ошибка обработки файла %s: %s", originalName, error),
	}

	err = s.Notify(
		detail,
		entities.User,
		guestFile.GuestID,
	)
	if err != nil {
		return fmt.Errorf("error save notification: %s", err.Error())
	}

	return nil
}
