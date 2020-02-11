package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	jwt "github.com/vizzyolog/showcase/jwt"
	yandexCLD "github.com/vizzyolog/showcase/yandexCLD"
)

const (
	botURL = "t.me/Anti_twaddle_bot"
)

func readTgTokenFromDisk() string {
	data, err := ioutil.ReadFile("telegramtoken.pem")
	if err != nil {
		fmt.Printf("\n err %v \n", err)
	}

	return string(data)
}

func downloadFile(url string, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	newIAM := jwt.GetNewIAMToken()
	tgToken := readTgTokenFromDisk()

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			// ignore any non-Message Updates
			continue
		}

		switch {
		case update.Message.Voice != nil:

			go recognize(update.Message, bot, tgToken, newIAM)

		case update.Message.Text != "":

			go synthesize(update.Message, bot, newIAM)

		}

	}

}

func synthesize(msg *tgbotapi.Message, bot *tgbotapi.BotAPI, newIAM string) {
	bytes, err := yandexCLD.SynthesizeVoice(newIAM, msg.Text)
	if err != nil {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ReplyToMessageID = msg.MessageID
		bot.Send(newMsg)
	}

	newFile := tgbotapi.FileBytes{
		Name:  "newVoiceFromYandex",
		Bytes: bytes,
	}

	newVoice := tgbotapi.NewVoiceUpload(msg.Chat.ID, newFile)
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "text")
	newMsg.ReplyMarkup = newVoice
	_, err = bot.Send(newVoice)
	if err != nil {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ReplyToMessageID = msg.MessageID
		bot.Send(newMsg)
	}

}

func recognize(msg *tgbotapi.Message, bot *tgbotapi.BotAPI, tgToken string, newIAM string) {
	voiceFileConfig := tgbotapi.FileConfig{
		FileID: msg.Voice.FileID,
	}
	file, err := bot.GetFile(voiceFileConfig)
	if err != nil {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ReplyToMessageID = msg.MessageID
		bot.Send(newMsg)
	}

	linkfordownload := file.Link(tgToken)

	downloadFile(linkfordownload, "tmp"+strconv.Itoa(msg.MessageID))

	voice := yandexCLD.ReadAudioFile("tmp" + strconv.Itoa(msg.MessageID))
	recoginzedText, err := yandexCLD.RecognizeVoice(newIAM, voice)
	if err != nil {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ReplyToMessageID = msg.MessageID
		bot.Send(newMsg)

		return
	}
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, recoginzedText)
	newMsg.ReplyToMessageID = msg.MessageID
	bot.Send(newMsg)
}
