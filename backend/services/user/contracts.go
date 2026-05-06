package user

type UserService interface {
	IsAuthenticated() bool
	UserId() (uint, error)
	InitGuestID() (uint, error)
	GuestID() (uint, error)
	UserPersonalDir() (string, error)
	GuestPersonalDir(isAbsolute bool) (string, error)
}
