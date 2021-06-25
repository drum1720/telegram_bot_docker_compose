package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var telegramReplyMessage TelegramReplyMessage

func main() {
	var settings Settings
	settings.updateData()
	offset := 0
	for {
		resp, err := http.Get(settings.BotUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
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
			workDirector(restResponse.Result[i], settings.BotUrl)
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
		task := TaskPhotoDownloader{
			ChatId: restResponse.Message.Chat.ChatId,
			Text:   restResponse.Message.Text,
			BotUrl: botUrl}
		task.sendTask()
	case "resize":
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
			FileId:   restResponse.Message.Photos[len(restResponse.Message.Photos)-1].FileId,
			FilePath: restResponse.Message.Photos[len(restResponse.Message.Photos)-1].FilePath,
		}
		task.sendTask()
	case "рифма":
		task := TaskRifma{
			ChatId: restResponse.Message.Chat.ChatId,
			Text:   restResponse.Message.Text,
			BotUrl: botUrl}
		task.sendTask()
	case "тор":
		task := TaskTorSearch{
			ChatId: restResponse.Message.Chat.ChatId,
			Text:   restResponse.Message.Text,
			BotUrl: botUrl}
		task.sendTask()
	case "help":
		replyMessage.Text = "напишите: 'photo + пробел + текст для поиска фото' для того чтобы получить рандомную фотографию по запросу " +
			"\n" + "чтобы уменьшить изображение, прикрепите его к сообщению, " +
			"само сообщение должно быть: 'resize + пробел + 1000 + 88', где первый параметр - размер по большей стороне, второй - степень сжатия(от 1 до 100), " +
			"второй параметр необязательный" +
			"\n" + "напишите: 'рифма + пробел + слово для поиска рифмы' чтобы получить все возможные рифмы к слову" +
			"\n" + "напишите: 'тор + пробел + текст для поиска торрент ссылок' чтобы получить 10 самых популярных результатов поиска"
		replyMessage.reply()
	default:
		replyMessage.Text = restResponse.Message.Chat.FirstName + ", " + "чет я не пойму, ты быканул(а) сейчас?"
		replyMessage.reply()
	}
}
