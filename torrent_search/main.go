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
	router.HandleFunc("/api/gettorrent", GetTorrent)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func GetTorrent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("body read all err", err)
		return
	}
	var taskTorSearch TaskTorSearch
	err = json.Unmarshal(body, &taskTorSearch)
	if err != nil {
		fmt.Println("Unmarshal BAD", err)
		return
	}
	textSearch := TrimTextSearch(taskTorSearch.Text, taskTorSearch.ChatId)
	if textSearch == "" {
		return
	}

	result := parseRifma(textSearch)
	if result == "" {
		result = "Такого нет даже на трекерах, долбанн-(ая/ый) ты извращен-(ец/ка)"
	}
	replyMessage := ReplyMessage{
		ChatId:                taskTorSearch.ChatId,
		Text:                  result,
		DisableWebPagePreview: true,
	}
	replyMessage.reply()
}

func parseRifma(textSearch string) string {
	var setting Settings
	setting.updateData()

	r, _ := http.NewRequest("GET", setting.TorSearchUrl+textSearch, nil)
	r.Header.Add("User-Agent", "PashaMan")
	r.Header.Add("Host", "55")
	client := &http.Client{}
	req, err := client.Do(r)
	if err != nil {
		fmt.Println("http err", err)
		return ""
	}
	defer req.Body.Close()

	bodyHtml, err := ioutil.ReadAll(req.Body)
	textHtml := string(bodyHtml)
	resultTrim := strings.Split(textHtml, "<p><a rel=\"nofollow\" target=\"_blank\"")
	var torSearchResults []StructsString
	for i, j := 1, 0; i < len(resultTrim)-1 && j < 10; i++ {
		if strings.Contains(resultTrim[i], "rutracker.org") {
			nameIurl := strings.Split(resultTrim[i], "\">")
			url := nameIurl[0]
			url = strings.ReplaceAll(url, "href=\"", "")
			url = strings.TrimSpace(url)
			name := nameIurl[1]
			name = strings.ReplaceAll(name, "<b>", "")
			name = strings.ReplaceAll(name, "</b>", "")
			name = strings.ReplaceAll(name, "<div class=\"h2", "")
			name = strings.ReplaceAll(name, "</a></p>", "")
			name = strings.TrimSpace(name)
			span := strings.Split(resultTrim[i], "span")
			size := span[2]
			size = strings.ReplaceAll(size, "class=\"size\">", "")
			size = strings.ReplaceAll(size, "</", "")
			size = strings.ReplaceAll(size, "&nbsp;", " ")
			seed := span[7]
			seed = strings.ReplaceAll(seed, "class=\"seeders\"><i class=\"fas fa-long-arrow-alt-up\"></i>", "")
			seed = strings.ReplaceAll(seed, "</", "")

			torSearchResult := TorSearchResult{
				TorTrackerName: "Rutracker",
				Name:           name,
				Url:            url,
				Size:           size,
				Seed:           seed,
			}
			torSearchResults = append(torSearchResults, torSearchResult)
			j++
		}
	}
	result := StructsToString(torSearchResults)
	return result
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
		if len(textSearch)+len(textSearchs[i]) < 250 {
			textSearch = textSearch + "%20" + textSearchs[i]
			continue
		}
		break
	}
	return textSearch
}
