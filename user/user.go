package user

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Secret    string
	AccessKey string
}
