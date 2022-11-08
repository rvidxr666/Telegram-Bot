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

func downloadImage(link string) *bytes.Buffer {
	res, err := http.Get(link)

	if err != nil {
		log.Fatalln(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	buffer := bytes.NewBuffer(data)

	return buffer
}

func sendImageToIBM(buf *bytes.Buffer) *http.Response {
	client := http.Client{}
	req, err := http.NewRequest("POST", url+"/v1/recognize", buf)

	if err != nil {
		log.Fatalln(err)
	}

	combination := "apiKey:" + os.Getenv("IBMWatsonKey")
	auth_header := baseEncoder(combination)

	req.Header.Set("Authorization", "Basic "+auth_header)
	req.Header.Set("Content-Type", "audio/ogg")

	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	return res
}

func parseResponse(res *http.Response) string {
	resultMap := make(map[string]interface{})
	data, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(data, &resultMap)

	if len(resultMap["results"].([]interface{})) == 0 {
		return "None"
	}

	text := resultMap["results"].([]interface{})[0].(map[string]interface{})["alternatives"].([]interface{})[0].(map[string]interface{})["transcript"]
	return text.(string)
}

func GetText(voice *tgbotapi.Voice) string {
	link := getLink(voice)
	buf := downloadImage(link)
	res := sendImageToIBM(buf)
	text := parseResponse(res)
	return text
}
