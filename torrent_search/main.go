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
	router.HandleFunc("/api/gettorrent", GetTorrent)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func GetTorrent(w http.ResponseWriter, r *http.Request) {
	var taskTorSearch TaskTorSearch
	taskTorSearch.UnmarshalBodyJson(r)
	telegramReplyMessage.ChatId = taskTorSearch.ChatId

	textSearch := TrimTextSearch(taskTorSearch.Text)
	if textSearch == "" {
		return
	}

	torSearchResults := getTorSearchResults(textSearch)
	if torSearchResults == "" {
		torSearchResults = "Такого нет даже на трекерах, долбанн-(ая/ый) ты извращен-(ец/ка)"
	}
	telegramReplyMessage.reply(torSearchResults)
}

func getTorSearchResults(textSearch string) string {
	var setting Settings
	setting.updateData()

	r, _ := http.NewRequest("GET", setting.TorSearchUrl+textSearch, nil)
	r.Header.Add("User-Agent", "ManInBlack")
	r.Header.Add("Host", "55")
	client := &http.Client{}
	req, err := client.Do(r)
	if err != nil {
		log.Println("http err", err)
		return ""
	}
	defer req.Body.Close()

	bodyHtml, err := ioutil.ReadAll(req.Body)
	torSearchResults := parseTorSearchResults(string(bodyHtml))
	result := TorStructsToString(torSearchResults)

	return result
}

func TrimTextSearch(Text string) string {
	textSearchs := strings.Split(Text, " ")
	if len(textSearchs) < 2 {
		telegramReplyMessage.reply("Как же я буду искать без текста для запроса...Заявка отклонена, возвращайтесь, как поймете что искать")
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

func parseTorSearchResults(textHtml string) []TorSearchResult {
	torrentTrackerName := "rutracker.org"
	resultCount := 10
	var torSearchResults []TorSearchResult

	resultTrim := strings.Split(textHtml, "<p><a rel=\"nofollow\" target=\"_blank\"")
	for i, j := 1, 0; i < len(resultTrim)-1 && j < resultCount; i++ {
		if strings.Contains(resultTrim[i], torrentTrackerName) {
			torSearchResult := parseTorSearchResult(resultTrim[i], torrentTrackerName)
			torSearchResults = append(torSearchResults, torSearchResult)
			j++
		}
	}
	return torSearchResults
}

func parseTorSearchResult(blockHtml string, torrentTrackerName string) TorSearchResult {
	nameIurl := strings.Split(blockHtml, "\">")

	url := nameIurl[0]
	url = strings.ReplaceAll(url, "href=\"", "")
	url = strings.TrimSpace(url)

	name := nameIurl[1]
	name = strings.ReplaceAll(name, "<b>", "")
	name = strings.ReplaceAll(name, "</b>", "")
	name = strings.ReplaceAll(name, "<div class=\"h2", "")
	name = strings.ReplaceAll(name, "</a></p>", "")
	name = strings.TrimSpace(name)

	sizeAndSeed := strings.Split(blockHtml, "span")

	size := sizeAndSeed[2]
	size = strings.ReplaceAll(size, "class=\"size\">", "")
	size = strings.ReplaceAll(size, "</", "")
	size = strings.ReplaceAll(size, "&nbsp;", " ")

	seed := sizeAndSeed[7]
	seed = strings.ReplaceAll(seed, "class=\"seeders\"><i class=\"fas fa-long-arrow-alt-up\"></i>", "")
	seed = strings.ReplaceAll(seed, "</", "")

	torSearchResult := TorSearchResult{
		TorTrackerName: torrentTrackerName,
		Name:           name,
		Url:            url,
		Size:           size,
		Seed:           seed,
	}
	return torSearchResult
}
