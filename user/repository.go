package user

type Reader interface {
	UserExistByEmail(email string) (bool, error)
	UserExistByKeyAndSecret(key string, secret string) (bool, error)
	FindUserByKeyAndSecret(key string, secret string) (*User, error)
	Login(user *User) (string, error)
}

type Writer interface {
	Create(user *User) error
	Delete(email string) error
}

//Repository repository interface
type UserRepository interface {
	Reader
	Writer
}