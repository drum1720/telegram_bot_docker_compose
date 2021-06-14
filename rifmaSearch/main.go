package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("полученный запрос не читается", err)
		return
	}
	var taskRifma TaskRifma
	err = json.Unmarshal(body, &taskRifma)
	if err != nil {
		fmt.Println("Унмаршал Бэд", err)
		return
	}
	textSearch := TrimTextSearch(taskRifma.Text, taskRifma.ChatId)
	if textSearch == "" {
		return
	}

	var db Db
	err = db.connect()
	if err != nil {
		replyMessage := ReplyMessage{ChatId: taskRifma.ChatId,
			Text: "База данных сломалась, попробуйте позже"}
		replyMessage.reply()
		return
	}
	var rifma = Rifma{Request: textSearch}
	err = rifma.WhereOneResponse(db)
	if err != nil {
		fmt.Println(err)
	}
	if rifma.Rifma != "" {
		replyMessage := ReplyMessage{ChatId: taskRifma.ChatId,
			Text: rifma.Rifma}
		replyMessage.reply()
		return
	}

	result:=parseRifma(textSearch)

	rifma.Rifma = result
	replyMessage := ReplyMessage{ChatId: taskRifma.ChatId,
		Text: result}
	replyMessage.reply()

	if rifma.Rifma != "" {
		rifma.AddToTable(db)
		return
	}
}

func parseRifma(textSearch string) string {
	req, err := http.Get("https://rifmus.net/rifma/" + textSearch)
	if err != nil {
		fmt.Println("данные не загружены", err)
		return ""
	}
	defer req.Body.Close()

	bodyHtml, err := ioutil.ReadAll(req.Body)
	textHtml := string(bodyHtml)
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
	return result
}

func TrimTextSearch(Text string, ChatId int) string {
	textSearchs := strings.Split(Text, " ")
	if len(textSearchs) < 2 {
		replyMessage := ReplyMessage{ChatId: ChatId,
			Text: "Не вижу слова я, для нахожденья рифмы, и я не слеп, да и глаза мои открыты..."}
		replyMessage.reply()
		return ""
	}
	return textSearchs[1]
}
