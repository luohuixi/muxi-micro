package example

import "context"

//go:generate mockgen -destination=mock_dao.go -package=example . UserDAO

type UserDAO interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

type User struct {
	ID   int64
	Name string
	Age  int
}

func (u *User) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return u, nil
}

func (u *User) CreateUser(ctx context.Context, user *User) (bool, error) {
	return true, nil
}
