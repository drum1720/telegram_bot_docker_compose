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
	PixabayImageDownloaderUrl string `json:"pixabay_image_downloader_url"`
	ImageHandlerResizeUrl     string `json:"image_handler_resize_url"`
	RifmaSearchUrl            string `json:"rifma_search_url"`
	TorrentSearchUrl          string `json:"torrent_search_url"`
	BotToken                  string `json:"bot_token"`
	BotApi                    string `json:"bot_api"`
	BotUrl                    string
}

type Chat struct {
	ChatId    int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Message struct {
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Photos   []Photo  `json:"photo"`
	Caption  string   `json:"caption"`
	Document Document `json:"document"`
}

type Document struct {
	FileId string `json:"file_id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type ReplyMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
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
