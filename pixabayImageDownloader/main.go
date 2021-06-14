package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	var settings Settings
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/getPhoto", DownloadPhoto)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func DownloadPhoto(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var taskPhoto TaskPhotoDownloader
	err = json.Unmarshal(body, &taskPhoto)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	textSearch := TrimTextSearch(taskPhoto.Text, taskPhoto.ChatId)
	if textSearch == "" {
		return
	}
	hit, err := getPixabayHit(textSearch)
	if hit.LargeImageURL == "" || err != nil {
		replyMessage := ReplyMessage{ChatId: taskPhoto.ChatId,
			Text: "Я пока умею искать только на обном сайте с картинками, и там, к сожалению такого нет((("}
		replyMessage.reply()
		return
	}
	replyMessage := ReplyMessage{ChatId: taskPhoto.ChatId,
		Text: hit.LargeImageURL}
	replyMessage.reply()
}

func getPixabayHit(textSearch string) (PixabayHit, error) {
	var settings Settings
	settings.updateData()
	resp, err := http.Get("https://pixabay.com/api/?key=" + settings.PixabayApiKey + "&q=" + textSearch + settings.PixabayParams)
	if err != nil {
		return PixabayHit{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err, "данные не читаются")
		return PixabayHit{}, err
	}
	var pixabayResponse PixabayResponse //структура для ответа от апи сервиса картинок
	err = json.Unmarshal(body, &pixabayResponse)
	if err != nil {
		fmt.Println(err, "JSON данные не читаются")
		return PixabayHit{}, err
	}
	if len(pixabayResponse.Hits) < 1 {
		return PixabayHit{}, err
	}
	if len(pixabayResponse.Hits) == 1 {
		return pixabayResponse.Hits[0], err
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(pixabayResponse.Hits) - 1)
	if index < 1 {
		for ; ; {
			index = rand.Intn(len(pixabayResponse.Hits) - 1)
			if index > 1 {
				break
			}
		}
	}
	return pixabayResponse.Hits[index], err
}

func TrimTextSearch(Text string, ChatId int) string {
	textSearchs := strings.Split(Text, " ")
	if len(textSearchs) < 2 {
		replyMessage := ReplyMessage{ChatId: ChatId,
			Text: "Как же я буду искать без текста для запроса...Заявка отклонена, возвращайтесь, как поймете что искать"}
		replyMessage.reply()
		return ""
	}
	textSearch := textSearchs[1]
	for i := 2; i < len(textSearchs); i++ {
		if len(textSearch)+len(textSearchs[i]) < 100 {
			textSearch = textSearch + "+" + textSearchs[i]
			continue
		}
		break
	}
	return textSearch
}
