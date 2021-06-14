package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	var settings Settings
	settings.updateData()
	botUrl := settings.BotUrl
	offset := 0
	for {

		resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
		if err != nil {
			fmt.Println("нет связи с апи телеграмм")
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		}
		var restResponse RestResponse
		err = json.Unmarshal(body, &restResponse)
		if err != nil {
		}
		for i := 0; i < len(restResponse.Result); i++ {
			workDirector(restResponse.Result[i], botUrl)
			offset = restResponse.Result[i].UpdateId + 1
		}
	}
}

func workDirector(restResponse Update, botUrl string) {
	var replyMessage ReplyMessage
	if restResponse.Message.Caption != "" {
		restResponse.Message.Text = restResponse.Message.Caption
	}
	text := strings.ToLower(restResponse.Message.Text)
	replyMessage.ChatId = restResponse.Message.Chat.ChatId
	message := strings.Split(text, " ")
	switch message[0] {
	case "photo":
		replyMessage.Text = "Привет " + restResponse.Message.Chat.FirstName + "! " + "Уже бегу искать фоточку"
		replyMessage.reply()
		task := TaskPhotoDownloader{
			ChatId: restResponse.Message.Chat.ChatId,
			Text:   restResponse.Message.Text,
			BotUrl: botUrl}
		task.sendTask()
	case "resize":
		replyMessage.Text = "Привет " + restResponse.Message.Chat.FirstName + "! " + "Сейчас сделаем)"
		replyMessage.reply()
		if restResponse.Message.Photos == nil {
			restResponse.Message.Photos = append(restResponse.Message.Photos, Photo{FileId: restResponse.Message.Document.FileId})
		}
		if restResponse.Message.Photos == nil {
			replyMessage.Text = "Где картинка??? "
			replyMessage.reply()
			return
		}
		restResponse.Message.Photos[len(restResponse.Message.Photos)-1].GetFileResult()
		task := TaskImageHandler{
			ChatId:   restResponse.Message.Chat.ChatId,
			Task:     text,
			BotUrl:   botUrl,
			FileId:   restResponse.Message.Photos[len(restResponse.Message.Photos)-1].FileId,
			FilePath: restResponse.Message.Photos[len(restResponse.Message.Photos)-1].FilePath,
		}
		task.sendTask()
	case "рифма":
		replyMessage.Text = "Привет " + restResponse.Message.Chat.FirstName + "! " + "Я всё устрою, не пройдет и пары лет. Иди готовь омлет)"
		replyMessage.reply()
		task := TaskRifma{
			ChatId: restResponse.Message.Chat.ChatId,
			Text:   restResponse.Message.Text,
			BotUrl: botUrl}
		task.sendTask()
	case "help":
		replyMessage.Text = "напишите: 'photo + пробел + текст для поиска фото' для того чтобы получить рандомную фотографию по запросу "
		replyMessage.reply()
		replyMessage.Text = "чтобы уменьшить изображение, прикрепите его к сообщению, " +
			"само сообщение должно быть: 'resize + пробел + 1000 + 88', где первый параметр - размер по большей стороне, второй - степень сжатия(от 1 до 100), " +
			"второй параметр необязательный"
		replyMessage.reply()
		replyMessage.Text = "напишите: 'рифма + пробел + слово для поиска рифмы' чтобы получить все возможные рифмы к слову"
		replyMessage.reply()
	default:
		replyMessage.Text = restResponse.Message.Chat.FirstName + ", " + "чет я не пойму, ты быканул(а) сейчас?"
		replyMessage.reply()
	}
}