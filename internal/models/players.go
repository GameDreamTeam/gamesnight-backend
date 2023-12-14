package models

type Player struct {
	Name             *string `json:"name"`
	Id               *string `json:"id"`
	PhrasesSubmitted bool    `json:"wordsSubmitted"`
}

type PlayerName struct {
	Name string `json:"name"`
}

type PlayerWords struct {
	Id         *string     `json:"id"`
	PhraseList *PhraseList `json:"phraseList"`
}
