package models

import "time"

type GameState int

const (
	PlayersJoining GameState = iota // 0
	AddingWords                     // 1
	Playing                         // 2
	Finished
)

type GameMeta struct {
	GameId    string    `json:"gameId"`
	AdminId   string    `json:"adminId"`
	Players   *[]Player `json:"players"`
	CreatedAt time.Time `json:"cAt"`
}

type Phrase struct {
	Input     string    `json:"input"`
	CreatedAt time.Time `json:"cAt"`
}

type UserInput struct {
	UserId  string    `json:"userId"`
	Phrases *[]Phrase `json:"phrases"`
}

type Game struct {
	GameId        string    `json:"gameId"`
	GameState     GameState `json:"state"`
	Teams         []Team    `json:"teams"`
	CurrentPlayer *Player   `json:"currentPlayer"`
}

type Team struct {
	Name    string    `json:"name"`
	Players *[]Player `json:"players"`
	Score   int       `json:"score"`
}
