package models

import "time"

type Game struct {
	GameId        string    `json:"gameId"`
	Admin         *Player   `json:"admin"`
	PlayerIds     *[]Player `json:"playerIds"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"cAt"`
	CurrentPlayer *Player   `json:"currentPlayer"`
}

type Phrase struct {
	Input     string    `json:"input"`
	CreatedAt time.Time `json:"cAt"`
}

type UserInput struct {
	UserId  string    `json:"userId"`
	Phrases *[]Phrase `json:"phrases"`
}
