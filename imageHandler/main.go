package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var telegramReplyMessage TelegramReplyMessage

func main() {
	var settings Settings
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/resize", Resize)
	http.Handle("/", router)
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func Resize(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var taskImageHandler TaskImageHandler
	taskImageHandler.Unmarshal(body)
	telegramReplyMessage.ChatId = taskImageHandler.ChatId

	size, quality := ExtractSizeAndQualityFromTask(taskImageHandler.Task, taskImageHandler.ChatId)
	modifiedImage, err := DownloadAndResizeImage(taskImageHandler, size, quality)
	if err != nil {
		log.Println(err)
		telegramReplyMessage.reply("Сожалею, что-то пошло не по плану")
		return
	}

	err = SendMultipartPostRequest(taskImageHandler, modifiedImage)
	if err != nil {
		log.Println(err)
		telegramReplyMessage.reply("Сожалею, что-то пошло не по плану")
		return
	}
}

func ExtractSizeAndQualityFromTask(Text string, ChatId int) (int, int) {
	size := 0
	quality := 100

	params := strings.Split(Text, " ")
	if len(params) < 2 {
		telegramReplyMessage.reply("Укажите требуемые размеры")
		return 0, 100
	}
	if len(params) > 2 {
		quality64, err := strconv.ParseInt(params[2], 10, 64)
		if err == nil && quality64 > 0 && quality64 <= 100 {
			quality = int(quality64)
		}
	}

	size64, err := strconv.ParseInt(params[1], 10, 64)
	if err == nil && size64 > 0 {
		size = int(size64)
	}

	return size, quality
}

func DownloadAndResizeImage(taskImageHandler TaskImageHandler, size int, quality int) (string, error) {
	resizeImageFullName := "reserv.jpg"

	resp, err := http.Get(taskImageHandler.FilePath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return "", err
	}
	if size > 4000 {
		size = 4000
	}
	resizeImg := resize.Thumbnail(uint(size), uint(size), img, resize.Lanczos3)

	out, err := os.Create(resizeImageFullName)
	if err != nil {
		return "", err //если не получилось создать файл, возвращаем ошибку
	}
	defer out.Close()

	err = jpeg.Encode(out, resizeImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return "", err
	}

	return resizeImageFullName, err
}

func SendMultipartPostRequest(taskImageHandler TaskImageHandler, filepath string) error {
	var settings Settings
	settings.updateData()

	file, _ := os.Open(filepath)
	defer file.Close()

	multipartBody := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartBody)
	part, _ := writer.CreateFormFile("document", file.Name())
	io.Copy(part, file)

	fieldWriter, err := writer.CreateFormField("chat_id")
	if err != nil {
		return err
	}
	_, err = fieldWriter.Write([]byte(strconv.Itoa(taskImageHandler.ChatId)))
	if err != nil {
		return err
	}
	writer.Close()

	r, _ := http.NewRequest("POST", settings.BotUrl+"/sendDocument", multipartBody)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	_, err = client.Do(r)

	return err
}
