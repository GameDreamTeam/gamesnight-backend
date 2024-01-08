package models

import "time"

type GameState int

const (
	PlayersJoining GameState = iota
	AddingWords
	TeamsDivided
	Playing
	Finished
)

type GameMeta struct {
	GameId    string    `json:"gameId"`
	AdminId   string    `json:"adminId"`
	Players   *[]Player `json:"players"`
	CreatedAt time.Time `json:"cAt"`
}

type Game struct {
	GameId                string    `json:"gameId"`
	GameState             GameState `json:"state"`
	Teams                 *[]Team   `json:"teams"`
	CurrentPlayer         *Player   `json:"currentPlayer"`
	NextPlayer            *Player   `json:"nextPlayer"`
	CurrentTeamIndex      int       `json:"currentTeamIndex"`
	CurrentPhraseMapIndex int       `json:"currentPhraseMapIndex"`
}

type Team struct {
	Name               string    `json:"name"`
	Players            *[]Player `json:"players"`
	Score              int       `json:"score"`
	CurrentPlayerIndex int       `json:"currentPlayerIndex"`
}
