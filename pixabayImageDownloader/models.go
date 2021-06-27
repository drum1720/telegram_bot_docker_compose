package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Settings struct {
	BotToken      string `json:"bot_token"`
	BotApi        string `json:"bot_api"`
	BotUrl        string
	PixabayApiKey string `json:"pixabay_api_key"`
	PixabayParams string `json:"pixabay_params"`
	ServerPort    string `json:"server_port"`
}

type TaskPhotoDownloader struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

type PixabayResponse struct {
	Hits []PixabayHit `json:"hits"`
}

type PixabayHit struct {
	LargeImageURL string `json:"largeImageURL"`
}

type TelegramReplyMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func (taskStruct *TaskPhotoDownloader) UnmarshalBodyJson(r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		log.Println(err)
		return
	}
}

func (reply *TelegramReplyMessage) reply(textMessage string) {
	var settings Settings
	settings.updateData()
	reply.Text = textMessage
	buf, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
	}
	_, err = http.Post(settings.BotUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Println(err, "не удалось отправить сообщение")
	}
}

func (s *Settings) updateData() {
	buff, err := ioutil.ReadFile("settings.json")
	if err != nil {
		panic("Файл настроек не читается")
	}
	err = json.Unmarshal(buff, &s)
	if err != nil {
		panic("Файл настроек не читается")
	}
	s.BotUrl = s.BotApi + s.BotToken
}
