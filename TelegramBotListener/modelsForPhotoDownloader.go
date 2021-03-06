package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type TaskPhotoDownloader struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

func (t TaskPhotoDownloader) sendTask() {
	var settings Settings
	settings.updateData()

	buf, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = http.Post(settings.PixabayImageDownloaderUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Println(err)
		return
	}
}
