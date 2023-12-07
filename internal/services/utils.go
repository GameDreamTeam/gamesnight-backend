package services

import (
	"errors"
	"gamesnight/internal/database"
	"gamesnight/internal/models"
	"math/rand"
	"time"
)

func contains(playerSlice []models.Player, player *models.Player) bool {
	for _, p := range playerSlice {
		if *p.Id == *player.Id {
			return true
		}
	}
	return false
}

func getNextTeamIndex(currentIndex int) int {
	return currentIndex ^ 1
}

func addPlayerToGame(gameMeta *models.GameMeta, player *models.Player) (*models.GameMeta, error) {

	if contains(*gameMeta.Players, player) {
		*gameMeta.Players = append(*gameMeta.Players, *player)
	} else {
		// Return custom error here (404)
		return nil, errors.New("player already exists in this game")
	}

	return gameMeta, nil
}

func dividePlayersIntoTeams(players []models.Player) ([]models.Player, []models.Player) {
	// if team exits in
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	mid := len(players) / 2
	return players[:mid], players[mid:]
}

func GeneratePhraseListToMap(phrases *models.PhraseList) (models.PhraseStatusMap, error) {
	if phrases == nil || phrases.List == nil {
		// return error
		return models.PhraseStatusMap{}, errors.New("no phrases found")
	}

	phraseStatusMap := models.PhraseStatusMap{
		Phrases: make([]models.Phrase, len(*phrases.List)),
		Status:  make([]models.PhraseStatus, len(*phrases.List)),
	}

	for i, phrase := range *phrases.List {
		phraseStatusMap.Phrases[i] = phrase
		phraseStatusMap.Status[i] = models.NotGuessed
	}

	return phraseStatusMap, nil
}

func (gs *GameService) StartTurnTimer(gameId string) error {
	game, err := database.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.GameState != models.Playing {
		game.GameState = models.Playing
		err = database.SetGame(game)
		if err != nil {
			return err
		}
	}

	turnDuration := 60 * time.Second
	timer := time.NewTimer(turnDuration)

	go func() {
		<-timer.C // This blocks until the timer expires

		//	 err := gs.function(gameId) this function should lead to displaying the changed currentphrases as a form
		// if err != nil {
		// 	logger.GetLogger().Logger.Error(
		// 		"error handling timer interrupt",
		// 		zap.Any("game", game),
		// 	)
		// }
	}()

	// Notify clients about the start of the turn
	// You can use a WebSocket or another communication mechanism for real-time updates

	return nil
}
