package models

import "time"

type Game struct {
	GameId    string    `json:"gameId"`
	AdminId   string    `json:"adminId"`
	CreatedAt time.Time `json:"cAt"`
}
