package services

import (
	"fmt"
	"gamesnight/internal/logger"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type SlackService struct{}

var ss *SlackService

func NewSlackService() {
	ss = &SlackService{}
}

func GetSlackService() *SlackService {
	return ss
}

// SendToSlack sends a message to Slack
func (ss *SlackService) SendToSlack(webhookURL, message string) error {
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, message)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response using Zap
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Logger.Error("Error reading Slack API response body", zap.Error(err))
		return err
	}

	logger.GetLogger().Logger.Info("Slack API Response",
		zap.String("ResponseBody", string(body)),
		zap.Int("StatusCode", resp.StatusCode),
	)

	return nil
}
