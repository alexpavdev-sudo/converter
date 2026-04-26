package user

import (
	"converter/app"
	"converter/entities"
	"converter/helpers"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"
)

const keyGuestId = "guest_id"
const keyUserId = "user_id"

type UserService struct {
	db      *gorm.DB
	session sessions.Session
}

func NewUserService(session sessions.Session) *UserService {
	return &UserService{
		db:      app.App().DB,
		session: session,
	}
}

func (s *UserService) generateGuestPersonalDir() (string, error) {
	for {
		randStr, err := helpers.GenerateRandomStoredName(16)
		if err != nil {
			return "", err
		}
		var exists bool
		err = s.db.Raw(
			"SELECT EXISTS(SELECT 1 FROM guests WHERE personal_dir = ?)",
			randStr,
		).Scan(&exists).Error
		if err != nil {
			return "", err
		}

		if !exists {
			return randStr, nil
		}
	}
}

func (s *UserService) IsAuthenticated() bool {
	_, err := s.UserId()
	if err == nil {
		return true
	}
	return false
}

func (s *UserService) UserId() (uint, error) {
	id := s.session.Get(keyUserId)
	if id != nil && id.(uint) > 0 {
		return id.(uint), nil
	}

	return 0, errors.New("the user is not authenticated")
}

func (s *UserService) GuestId() (uint, error) {
	id := s.session.Get(keyGuestId)
	if id != nil && id.(uint) > 0 {
		var guest entities.Guest
		if err := s.db.First(&guest, id.(uint)).Error; err == nil {
			return id.(uint), nil
		}
	}
	return s.initGuest()
}

func (s *UserService) initGuest() (uint, error) {
	personalDir, err := s.generateGuestPersonalDir()
	if err != nil {
		return 0, err
	}
	guest := &entities.Guest{
		PersonalDir: personalDir,
	}
	if err := s.db.Create(guest).Error; err != nil {
		return 0, errors.New(fmt.Sprintf("database save failed: %s", err.Error()))
	}

	s.session.Set(keyGuestId, guest.ID)
	if err := s.session.Save(); err != nil {
		return 0, errors.New("error save session guest id")
	}
	return guest.ID, nil
}

func (s *UserService) UserPersonalDir() (string, error) {
	//todo
	return "", errors.New("error")
}

func (s *UserService) GuestPersonalDir(isAbsolute bool) (string, error) {
	id, err := s.GuestId()
	if err != nil {
		return "", err
	}

	var guest entities.Guest
	if err := s.db.First(&guest, id).Error; err != nil {
		return "", err
	}

	return guest.PersonalPath(isAbsolute), nil
}
