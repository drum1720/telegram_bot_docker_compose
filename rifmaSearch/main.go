package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var telegramReplyMessage TelegramReplyMessage

func main() {
	var settings Settings
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/getrifma", GetRifma)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func GetRifma(w http.ResponseWriter, r *http.Request) {
	var taskRifma TaskRifma
	taskRifma.UnmarshalBodyJson(r)
	telegramReplyMessage.ChatId = taskRifma.ChatId

	textSearch := TrimTextSearch(taskRifma.Text)
	if textSearch == "" {
		return
	}

	var db Db
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
	fmt.Println("i parse")

	if rifma.Rifma != "" {
		rifma.AddToTable(db)
		return
	}
}

func parseRifma(textSearch string) string {
	req, err := http.Get("https://rifmus.net/rifma/" + textSearch)
	if err != nil {
		telegramReplyMessage.reply("Я не смог придумать(")
		log.Println("http err", err)
		return ""
	}
	defer req.Body.Close()
	bodyHtml, err := ioutil.ReadAll(req.Body)
	textHtml := string(bodyHtml)

	//да, я идиот, но парсить слова-рифмы мне было интересно именно так
	resultTrim := strings.Split(textHtml, "<ul class='multicolumn' itemprop='text'>")
	var resultTrim1 []string
	for i := 1; i < len(resultTrim); i++ {
		resultTrim1 = append(resultTrim1, strings.Split(resultTrim[i], "</ul>")[0])
	}
	for i := 0; i < len(resultTrim1); i++ {
		s := strings.ReplaceAll(resultTrim1[i], "<li>", "")
		resultTrim1[i] = s
		s = strings.ReplaceAll(resultTrim1[i], "</li>", "")
		resultTrim1[i] = s
		s = strings.ReplaceAll(resultTrim1[i], "\n", " + ")
		resultTrim1[i] = s
	}
	var result string
	for i := 0; i < len(resultTrim1); i++ {
		result = result + resultTrim1[i]
	}
	//тут извращения заканчиваются

	return result
}

func TrimTextSearch(Text string) string {
	KeywordsToSearch := strings.Split(Text, " ")
	if len(KeywordsToSearch) < 2 {
		telegramReplyMessage.reply("Не вижу слова я, для нахожденья рифмы, и я не слеп, да и глаза мои открыты...")
		return ""
	}
	return KeywordsToSearch[1]
}
