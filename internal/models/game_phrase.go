package models

type PhraseStatus int

const (
	NotGuessed PhraseStatus = iota
	Guessed
)

var (
	CurrentIndex int = -1
)

type Phrase struct {
	Input string `json:"input"`
}

type PhraseList struct {
	List *[]Phrase `json:"phraseList"`
}

type PhraseStatusMap struct {
	Phrases map[string]PhraseStatus `json:"phraseListMap"`
}

type PlayerGuess struct {
	PlayerChoice string `json:"playerChoice"`
}

type ResponseData struct {
	Game       *Game
	NextPhrase string `json:"nextPhrase"`
}
