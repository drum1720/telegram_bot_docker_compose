package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
			log.Println("нет связи с апи телеграмм")
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("notBody")
			continue
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
	if restResponse.Message.Caption != "" {
		restResponse.Message.Text = restResponse.Message.Caption
	}

	text := strings.ToLower(restResponse.Message.Text)
	telegramReplyMessage.ChatId = restResponse.Message.Chat.ChatId
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
			telegramReplyMessage.reply("Где картинка???")
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
		}
		task.sendTask()
	case "help":
		telegramReplyMessage.reply("Напишите: '**photo**' + пробел + текст для поиска фото' для того чтобы получить рандомную фотографию по запросу " +
			"\n" + "\n" + "Чтобы уменьшить изображение, прикрепите его к сообщению: '**resize** + пробел + 1000 + 88', где первый параметр - размер по большей стороне, второй - степень сжатия(от 1 до 100), " +
			"второй параметр необязательный" +
			"\n" + "\n" + "Напишите: '**рифма** + пробел + слово для поиска рифмы' чтобы получить все возможные рифмы к слову" +
			"\n" + "\n" + "Напишите: '**тор** + пробел + текст для поиска торрент ссылок' чтобы получить 10 самых популярных результатов поиска")
	default:
		telegramReplyMessage.reply(restResponse.Message.Chat.FirstName + ", " + " я не пойму, ты быканул(а) сейчас?")
	}
}
