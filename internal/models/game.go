package models

import "time"

type Game struct {
	GameId    string    `json:"gameId"`
	Admin     *Player   `json:"admin"`
	CreatedAt time.Time `json:"cAt"`
	PlayerIds *[]Player `json:"playerIds"`
}

type GameWords struct {
	GameId     string `json:"gameId"`
	PhraseList *PhraseList
}

type Phrase struct {
	Input string `json:"input"`
}

type PhraseList struct {
	List *[]Phrase `json:"phraseList"`
}
