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
	BotToken     string `json:"bot_token"`
	BotApi       string `json:"bot_api"`
	BotUrl       string
	ServerPort   string `json:"server_port"`
	TorSearchUrl string `json:"tor_search_url"`
}

type TaskTorSearch struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	BotUrl string `json:"botUrl"`
}

type TorSearchResult struct {
	Url            string
	Name           string
	TorTrackerName string
	Seed           string
	Size           string
}
type ReplyMessage struct {
	ChatId                int    `json:"chat_id"`
	Text                  string `json:"text"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func (tsr TorSearchResult) ToString() string {
	if tsr.Url != "" {
		newLine := "\n"
		return tsr.Url + newLine + tsr.Name + newLine + tsr.TorTrackerName + newLine + "Size: " + tsr.Size + newLine + "Seeds: " + tsr.Seed + newLine
	}
	return ""
}

type StructsString interface {
	ToString() string
}

func StructsToString(StructsStrings []StructsString) string {
	var result string
	for i := 0; i < len(StructsStrings); i++ {
		result = result + StructsStrings[i].ToString() + "\n"
	}
	return result
}

func (taskStruct *TaskImageHandler) Unmarshal(requestBody []byte) {
	err := json.Unmarshal(requestBody, &taskStruct)
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
