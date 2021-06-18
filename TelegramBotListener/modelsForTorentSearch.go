package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TaskTorSearch struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

func (t TaskTorSearch) sendTask() {
	var settings Settings
	settings.updateData()
	buf, err := json.Marshal(t)
	if err != nil {
	}
	_, err = http.Post(settings.TorrentSearchUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
	}
}
