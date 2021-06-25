package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Settings struct {
	BotToken   string `json:"bot_token"`
	BotApi     string `json:"bot_api"`
	BotUrl     string
	ServerPort string `json:"server_port"`
	PgHost     string `json:"pg_host"`
	PgPort     string `json:"pg_port"`
	PgUser     string `json:"pg_user"`
	PgPass     string `json:"pg_pass"`
	PgDbName   string `json:"pg_db_name"`
}

type TaskRifma struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

type TelegramReplyMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func (taskStruct *TaskRifma) UnmarshalBodyJson(r *http.Request) {
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
	s.BotUrl = s.BotApi + s.BotToken
}
