package models

import "time"

type GameMeta struct {
	GameId        string    `json:"gameId"`
	AdminId       string    `json:"adminId"`
	Players       *[]Player `json:"players"`
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
