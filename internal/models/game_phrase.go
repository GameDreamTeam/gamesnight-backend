package models

import "time"

type PhraseStatus int

const (
	NotGuessed PhraseStatus = iota
	Guessed
)

var (
	CurrentIndex  int = 0
	TurnStartTime time.Time
)

type Phrase struct {
	Input string `json:"input"`
}

type PhraseList struct {
	List *[]Phrase `json:"phraseList"`
}

type PhraseStatusMap struct {
	Phrases []Phrase       `json:"phraseList"`
	Status  []PhraseStatus `json:"statusList"`
}

type PlayerGuess struct {
	PlayerChoice string `json:"playerChoice"`
}

type ResponseData struct {
	PhraseMap     *PhraseStatusMap `json:"phraseListMap"`
	CurrentPhrase string           `json:"currentPhrase"`
	TurnStartedAt time.Time        `json:"turnStartedAt"`
}
