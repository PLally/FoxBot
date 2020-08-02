package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type DiscordLogger struct {
	LogLevels []logrus.Level
	webhookUrl string
	httpClient *http.Client
}

func (hook DiscordLogger) Levels() []logrus.Level {
	return hook.LogLevels
}

func (hook DiscordLogger) Fire(entry *logrus.Entry) error {
	message := fmt.Sprintf("```%s : %s\n%s```", entry.Level.String(), entry.Time.Format(time.ANSIC), entry.Message)
	content, _ := json.Marshal(map[string]string{
		"content": message,
	})

	reader := bytes.NewReader(content)
	req, _ := http.NewRequest("POST", hook.webhookUrl, reader)
	req.Header.Set("Content-Type", "application/json")

	if entry.Context != nil {
		req.WithContext(entry.Context)
	}
	_, _ = hook.httpClient.Do(req)

	return nil
}

func WebhookHook(webhook string) logrus.Hook {
	return DiscordLogger{
		LogLevels: []logrus.Level{
			logrus.ErrorLevel,
		},
		httpClient: &http.Client{},
		webhookUrl: webhook,
	}
}