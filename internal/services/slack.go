package services

import (
	"fmt"
	"gamesnight/internal/logger"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// SendToSlack sends a message to Slack
func SendToSlack(webhookURL, message string) error {
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, message)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response using Zap
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.GetLogger().Logger.Info("Slack API Response",
		zap.String("ResponseBody", string(body)),
		zap.Int("StatusCode", resp.StatusCode),
	)

	return nil
}
