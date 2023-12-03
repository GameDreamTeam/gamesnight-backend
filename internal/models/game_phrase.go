package models

type PhraseStatus int

const (
	NotGuessed PhraseStatus = iota
	Guessed
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
