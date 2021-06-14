package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Settings struct {
	BotToken   string  `json:"bot_token"`
	BotApi     string  `json:"bot_api"`
	BotUrl     string
	ServerPort    string `json:"server_port"`
}

type ReplyMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type TaskImageHandler struct {
	Task     string `json:"task"`
	ChatId   int    `json:"chat_id"`
	BotUrl   string `json:"botUrl"`
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

func (reply ReplyMessage) reply() {
	var settings Settings
	settings.updateData()
	buf, err := json.Marshal(reply)
	if err != nil {
		fmt.Println(err)
	}
	_, err = http.Post(settings.BotUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err, "не удалось отправить сообщение")
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
	s.BotUrl=s.BotApi+s.BotToken
}
