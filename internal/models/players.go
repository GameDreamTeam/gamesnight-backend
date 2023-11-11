package models

type Player struct {
	Name *string `json:"name"`
	Id   *string `json:"id"`
}

type PlayerName struct {
	Username string `json:"username"`
}
