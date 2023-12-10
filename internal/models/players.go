package models

type Player struct {
	Name *string `json:"name"`
	Id   *string `json:"id"`
}

type PlayerName struct {
	Name string `json:"name"`
}

type PlayerWords struct {
	Id         *string `json:"id"`
	PhraseList *PhraseList
}
