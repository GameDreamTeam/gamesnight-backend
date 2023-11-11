package models

import "time"

type Game struct {
	GameId    string    `json:"gameId"`
	Admin     *Player   `json:"admin"`
	CreatedAt time.Time `json:"cAt"`
	PlayerIds *[]Player `json:"playerIds"`
}
