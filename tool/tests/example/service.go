package example

import (
	"context"
	"errors"
)

type UserService struct {
	userDAO UserDAO
}

func NewUserService(userDAO UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.userDAO.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, name string, age int) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if age <= 0 {
		return errors.New("age must be positive")
	}

	return s.userDAO.CreateUser(ctx, &User{
		Name: name,
		Age:  age,
	})
}
