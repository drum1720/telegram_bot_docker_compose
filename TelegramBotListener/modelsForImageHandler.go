package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (t TaskImageHandler) sendTask() {
	var settings Settings
	settings.updateData()
	buf, err := json.Marshal(t)
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
	}
	fmt.Println(settings.BotUrl + "/getFile")
	resp, err := http.Post(settings.BotUrl+"/getFile", "application/json", bytes.NewBuffer(buf))
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	err = json.Unmarshal(body, &photoResult)
	p.FilePath = "https://api.telegram.org/file/bot" + settings.BotToken + "/" + photoResult.Photo.FilePath
}
