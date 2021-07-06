package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var telegramReplyMessage TelegramReplyMessage
var db Db

func main() {
	err := db.connect()
	if err != nil {
		return
	}

	go dbSlowFilling()

	var settings Settings
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/getrifma", GetRifma)
	http.Handle("/", router)
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func GetRifma(w http.ResponseWriter, r *http.Request) {
	var taskRifma TaskRifma
	taskRifma.UnmarshalBodyJson(r)
	telegramReplyMessage.ChatId = taskRifma.ChatId

	textSearch := trimTextSearch(taskRifma.Text)
	if textSearch == "" {
		return
	}

	err := db.connect()
	if err != nil {
		telegramReplyMessage.reply("База данных сломалась, попробуйте позже")
		return
	}
	var rifma = Rifma{Request: textSearch}
	err = rifma.WhereOneResponse(db)
	if err != nil {
		fmt.Println(err)
	}
	if rifma.Rifma != "" {
		telegramReplyMessage.reply(rifma.Rifma)
		fmt.Println("db search")
		return
	}

	result := parseRifma(textSearch)
	rifma.Rifma = result
	telegramReplyMessage.reply(result)
	fmt.Println("i parse", result)

	if rifma.Rifma != "" {
		rifma.AddToTable(db)
		return
	}
}

func parseRifma(textSearch string) string {
	req, err := http.Get("https://rifmu.ru/" + textSearch)
	if err != nil {
		log.Println("http err", err)
		return ""
	}
	defer req.Body.Close()
	bodyHtml, err := ioutil.ReadAll(req.Body)
	textHtml := string(bodyHtml)

	//да, я идиот, но парсить слова-рифмы мне было интересно именно так
	var result string
	if strings.Contains(textHtml, "<ul class=\"row\" style=\"list-style-type: none; font-size: 18px;\">") {
		resultTrim := strings.Split(textHtml, "<ul class=\"row\" style=\"list-style-type: none; font-size: 18px;\">")
		resultTrim = strings.Split(resultTrim[1], "<li class=\"col-lg-4\" style='margin: 5px 0;'><a href=\"")
		for i := 1; i < len(resultTrim); i++ {
			result = result + "+" + strings.Split(resultTrim[i], "\"")[0]
		}
	}
	//тут извращения заканчиваются

	return result
}

func trimTextSearch(Text string) string {
	KeywordsToSearch := strings.Split(Text, " ")
	if len(KeywordsToSearch) < 2 {
		telegramReplyMessage.reply("Не вижу слова я, для нахожденья рифмы, и я не слеп, да и глаза мои открыты...")
		return ""
	}
	return KeywordsToSearch[1]
}

func dbSlowFilling() {
	for {
		time.Sleep(180 * time.Second)

		var word Word
		word.WhereOneResponse(db)
		if word.Word == "" {
			break
		}

		parseResult := parseRifma(word.Word)
		if parseResult != "" {
			rifma := Rifma{
				Request: word.Word,
				Rifma:   parseResult}
			rifma.AddToTable(db)
			word.Status = "Done"

		} else {
			word.Status = "None"
		}

		word.UpdateStatus(db)
	}
}
