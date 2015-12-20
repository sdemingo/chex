package users

type AppUser interface {
	ID() int64
	GetInfo() map[string]string
	GetRole() int8
	GetEmail() string
}
