package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rvidxr666/test-telegram-bot/audio"
	"github.com/rvidxr666/test-telegram-bot/weather"
)

// Global variables
var bot, err = tgbotapi.NewBotAPI(os.Getenv("telegramAPI"))

func initialMessage(ID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(ID, "Hello! It's a first version of my bot!")
	return msg
}

func landingTable(ID int64) tgbotapi.MessageConfig {
	message := `
	Please choose the option:
	/weather - extract weather data in your location
	/voice   - transform voice message to text  
	`

	msg := tgbotapi.NewMessage(ID, message)
	btnWeather := tgbotapi.KeyboardButton{
		Text: "/weather",
	}
	btnExtractText := tgbotapi.KeyboardButton{
		Text: "/voice",
	}
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btnWeather, btnExtractText})
	msg.ReplyMarkup = keyboard
	return msg
}

func landingPage(ID int64) {
	// initalMsg := initialMessage(ID)
	// bot.Send(initalMsg)
	tableMsg := landingTable(ID)
	bot.Send(tableMsg)
}

func checkWeather(ID int64) {
	msg := tgbotapi.NewMessage(ID, "Please share your location")
	btn := tgbotapi.KeyboardButton{
		Text:            "Share Location",
		RequestLocation: true,
	}
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func weatherResponse(location *tgbotapi.Location, ID int64) {
	long := location.Longitude
	lat := location.Latitude
	weatherMap := weather.QueryWeatherApi(long, lat)
	msg := tgbotapi.NewMessage(ID, "Location: "+weatherMap["Location"].(string)+"\nWeather: "+weatherMap["Weather"].(string)+
		"\nTemperature: "+fmt.Sprintf("%f", weatherMap["Temp"].(float64))[:4]+"°")

	bot.Send(msg)
	landingPage(ID)
}

func checkVoice(ID int64) {
	msg := tgbotapi.NewMessage(ID, "Please send a Voice Message.\nNOTE: It should not be longer than 1 minute.")
	bot.Send(msg)
}

func voiceResponse(messaage string, ID int64) {
	msg := tgbotapi.NewMessage(ID, "Voice Message: "+messaage)
	bot.Send(msg)
	landingPage(ID)
}

func main() {
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		id := update.Message.Chat.ID

		fmt.Println("Printed message", update.Message.Location)
		switch update.Message.Text {
		case "/start":
			landingPage(id)
		case "/weather":
			checkWeather(id)
		case "/voice":
			checkVoice(id)
		}

		// Start of weather logic (Think of improvement)
		location := update.Message.Location
		if location != nil {
			weatherResponse(location, id)
		}

		// Start of audio logic
		voice := update.Message.Voice
		fmt.Println("Voice", voice)

		if voice != nil && voice.Duration <= 60 {
			fmt.Println("Inside")
			text := audio.GetText(voice)
			voiceResponse(text, id)
		}
	}
}
