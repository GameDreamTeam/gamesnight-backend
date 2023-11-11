package services

import (
	"fmt"
	"gamesnight/internal/models"
)

type UserService struct{}

var us *UserService

func NewUserService() {
	us = &UserService{}
}

func GetUserService() *UserService {
	return us
}

func (us *UserService) CreateNewUser() (*models.User, error) {
	key, err := GetKeyGenerator().CreateUserKey()

	if err != nil {
		fmt.Println("Error in creating user key %s", err)
		return nil, err
	}

	user := &models.User{
		UserId: &key,
	}
	return user, nil
}
