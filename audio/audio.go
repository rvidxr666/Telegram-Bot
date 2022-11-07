package audio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string = os.Getenv("telegramAPI") // Transform to env variable
var bot, _ = tgbotapi.NewBotAPI(token)
var url string = os.Getenv("IBMWatsonURL") //Transform to env var

func baseEncoder(combination string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(combination))
	return sEnc
}

func getLink(voice *tgbotapi.Voice) string {
	flConf := tgbotapi.FileConfig{FileID: voice.FileID}
	fileObject, _ := bot.GetFile(flConf)
	fmt.Println(fileObject)
	fmt.Println("LINK LINK LINK: ", fileObject.Link(token))
	link := fileObject.Link(token)
	return link
}

func makeRequest(link string) {
	res, err := http.Get(link)

	if err != nil {
		log.Fatalln(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	buffer := bytes.NewBuffer(data)

	client := http.Client{}
	req, err := http.NewRequest("POST", url+"/v1/recognize", buffer)

	if err != nil {
		log.Fatalln(err)
	}

	combination := "apiKey:" + os.Getenv("IBMWatsonKey")
	auth_header := baseEncoder(combination)

	req.Header.Set("Authorization", "Basic "+auth_header)
	req.Header.Set("Content-Type", "audio/ogg")
	res, err = client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	resultMap := make(map[string]interface{})
	data, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(data, &resultMap)

	fmt.Println(resultMap["results"].([]interface{})[0].(map[string]interface{})["alternatives"].([]interface{})[0].(map[string]interface{})["transcript"])
}

func downloadFile(voice *tgbotapi.Voice) {
	link := getLink(voice)
	makeRequest(link)
	// res, err := http.Get(link)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // reqMap := make(map[string]interface{})
	// data, _ := ioutil.ReadAll(res.Body)
	// buffer := bytes.NewBuffer(data)

	// client := http.Client{}
	// req, _ := http.NewRequest("POST", url+"/v1/recognize", buffer)
	// req.Header.Set("Authorization", "Bearer "+os.Getenv("IBMWatsonKety"))
	// req.Header.Set("Content-Type", "audio/flac")
	// res, err = client.Do(req)
	// fmt.Println(bytes.NewBuffer(data))
}

func GetText(voice *tgbotapi.Voice) {
	downloadFile(voice)
	// fileID := voice.FileID
	// FileUniqueID := voice.FileUniqueID
	// fileObject := tgbotapi.File{FileID: fileID, FileUniqueID: FileUniqueID}
	// fmt.Println("LINK LINK LINK: ", fileObject.Link(token))
}
