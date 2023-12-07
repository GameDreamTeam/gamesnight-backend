package services

type UtilService struct{}

var us *GameService

func NewUtilService() {
	us = &GameService{}
}

func GetUtilService() *GameService {
	return us
}
