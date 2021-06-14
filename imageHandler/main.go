package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var settings Settings
	settings.updateData()
	router := mux.NewRouter() //объявление переадресации
	router.HandleFunc("/api/resize", Resize)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(settings.ServerPort, nil) //запуск сервера
}

func Resize(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var taskPhoto TaskImageHandler
	err = json.Unmarshal(body, &taskPhoto)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	size, quality := TrimParams(taskPhoto.Task, taskPhoto.ChatId)
	fileDone,err:=DownloadAndResizeImage(taskPhoto, size, quality)
	if err != nil {
		replyMessage := ReplyMessage{
			ChatId: taskPhoto.ChatId,
			Text:   "Сожалею, что-то пошло не по плану"}
		replyMessage.reply()
		return
	}
	err=SendMultipartPostRequest(taskPhoto,fileDone)
	if err != nil {
		replyMessage := ReplyMessage{
			ChatId: taskPhoto.ChatId,
			Text:   "Сожалею, что-то пошло не по плану"}
		replyMessage.reply()
		return
	}
}

func TrimParams(Text string, ChatId int) (int, int) {
	size := 0
	quality := 100
	params := strings.Split(Text, " ")
	if len(params) < 2 {
		replyMessage := ReplyMessage{
			ChatId: ChatId,
			Text:   "Укажите требуемые размеры"}
		replyMessage.reply()
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

func DownloadAndResizeImage(taskPhoto TaskImageHandler, size int, quality int) (string,error) {
	resp, err := http.Get(taskPhoto.FilePath)
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
	imgResize := resize.Thumbnail(uint(size), uint(size), img, resize.Lanczos3)
	out, err := os.Create("reserv.jpg")
	if err != nil {
		return "", err //если не получилось создать файл, возвращаем ошибку
	}
	defer out.Close()
	err = jpeg.Encode(out, imgResize, &jpeg.Options{Quality: quality})
	if err != nil {
		return "", err
	}
	return "reserv.jpg", err
}

func SendMultipartPostRequest(taskPhoto TaskImageHandler,filepath string)error{
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
	_, err = fieldWriter.Write([]byte(strconv.Itoa(taskPhoto.ChatId)))
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
