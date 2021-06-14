package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TaskRifma struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

func (t TaskRifma) sendTask() {
	var settings Settings
	settings.updateData()
	buf, err := json.Marshal(t)
	if err != nil {
	}
	_, err = http.Post(settings.RifmaSearchUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
	}
}
