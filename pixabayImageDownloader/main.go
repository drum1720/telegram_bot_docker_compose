package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var telegramReplyMessage TelegramReplyMessage
var settings Settings

func main() {
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/getPhoto", getPhoto)
	http.Handle("/", router)
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	var taskPhotoDownloader TaskPhotoDownloader
	taskPhotoDownloader.UnmarshalBodyJson(r)
	telegramReplyMessage.ChatId = taskPhotoDownloader.ChatId

	textForSearch := prepareTextSearch(taskPhotoDownloader.Text)
	if textForSearch == "" {
		return
	}

	hit, err := getRandomPixabayHit(textForSearch)
	if hit.LargeImageURL == "" || err != nil {
		telegramReplyMessage.reply("Я пока умею искать только на обном сайте с картинками, и там, к сожалению такого нет(((")
	}

	telegramReplyMessage.reply(hit.LargeImageURL)
}

func getRandomPixabayHit(textForSearch string) (PixabayHit, error) {
	settings.updateData()

	resp, err := http.Get("https://pixabay.com/api/?key=" + settings.PixabayApiKey + "&q=" + textForSearch + settings.PixabayParams)
	if err != nil {
		return PixabayHit{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, "данные не читаются")
		return PixabayHit{}, err
	}
	var pixabayResponse PixabayResponse
	err = json.Unmarshal(body, &pixabayResponse)
	if err != nil {
		log.Println(err, "JSON данные не читаются")
		return PixabayHit{}, err
	}

	if len(pixabayResponse.Hits) < 1 {
		return PixabayHit{}, err
	}
	if len(pixabayResponse.Hits) == 1 {
		return pixabayResponse.Hits[0], nil
	}

	randIndex := getRandIndex(pixabayResponse.Hits)
	return pixabayResponse.Hits[randIndex], err
}

func prepareTextSearch(Text string) string {
	maxLenTextSearch := 100
	KeywordsToSearch := strings.Split(Text, " ")
	if len(KeywordsToSearch) < 2 {
		telegramReplyMessage.reply("Как же я буду искать без текста для запроса...Заявка отклонена, возвращайтесь, как поймете что искать")
		return ""
	}
	textForSearch := KeywordsToSearch[1]
	for i := 2; i < len(KeywordsToSearch); i++ {
		if len(textForSearch)+len(KeywordsToSearch[i]) < maxLenTextSearch {
			textForSearch = textForSearch + "+" + KeywordsToSearch[i]
			continue
		}
		break
	}
	return textForSearch
}

func getRandIndex(hits []PixabayHit) int {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(hits) - 1)
	if randIndex < 1 {
		for {
			randIndex = rand.Intn(len(hits) - 1)
			if randIndex > 1 {
				break
			}
		}
	}
	return randIndex
}
