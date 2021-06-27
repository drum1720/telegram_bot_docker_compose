package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type TaskImageHandler struct {
	Task     string `json:"task"`
	ChatId   int    `json:"chat_id"`
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type Photo struct {
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type PhotoResult struct {
	Photo Photo `json:"result"`
}

func (task TaskImageHandler) sendTask() {
	var settings Settings
	settings.updateData()
	buf, err := json.Marshal(task)
	if err != nil {
	}
	_, err = http.Post(settings.ImageHandlerResizeUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
	}
}

func (p *Photo) GetFileResult() {
	var photoResult PhotoResult
	var settings Settings

	settings.updateData()
	buf, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(settings.BotUrl+"/getFile", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &photoResult)
	if err != nil {
		log.Println(err)
		return
	}

	p.FilePath = "https://api.telegram.org/file/bot" + settings.BotToken + "/" + photoResult.Photo.FilePath
}
