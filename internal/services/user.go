package services

import (
	"fmt"
	"gamesnight/internal/models"

	"github.com/gin-gonic/gin"
)

type UserService struct{}

var us *UserService

func NewUserService() {
	us = &UserService{}
}

func GetUserService() *UserService {
	return us
}

func CreateNewUser() (*models.User, error) {
	key, err := GetKeyGenerator().CreateUserKey()

	if err != nil {
		fmt.Printf("Error in creating user key %s", err)
		return nil, err
	}

	user := &models.User{
		UserId: &key,
	}
	return user, nil
}

func (us *UserService) GetUser(c *gin.Context, userCookieName string) (*models.User, error) {

	userCookie, err := c.Cookie(userCookieName)
	if err != nil {
		user, err := CreateNewUser()

		if err != nil {
			fmt.Printf("Error in creating new user %s", err)
			return nil, err
		}

		token, err := GetTokenService().CreateUserToken(*user.UserId)
		if err != nil {
			fmt.Printf("Error in creating new token %s", err)
			return nil, err
		}
		// Need to check these other parameters
		c.SetCookie(userCookieName, token.Token, 3600, "/", "", false, true)
		userCookie = token.Token
	}

	user, err := GetTokenService().ParseUserToken(userCookie)
	if err != nil {
		// Create new user in this case
		// Check if expired or not able to decode this token
		fmt.Printf("Error in parsing token %s", err)
		return nil, err
	}

	return user, nil
}
