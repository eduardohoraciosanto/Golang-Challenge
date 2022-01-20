package code_implementation

type UserRepository interface {
	GetUser(id string) (User, error)
}

type User struct {
	Id      string
	Friends []string
}
