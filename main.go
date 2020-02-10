package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	// voice := yandexCLD.ReadAudioFile("sources/speech3.ogg")
	// yandexCLD.SendPost(newIAM, voice)

	tgToken := readTgTokenFromDisk()
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Voice != nil {
			inVoice := update.Message.Voice
			fileConfig := tgbotapi.FileConfig{
				FileID: inVoice.FileID,
			}
			file, err := bot.GetFile(fileConfig)
			if err != nil {
				fmt.Println("Can't download the file", err)
			}

			linkfordownload := file.Link(tgToken)
			downloadFile(linkfordownload, "tmp.ogg")

			voice := yandexCLD.ReadAudioFile("tmp.ogg")
			recoginzedText := yandexCLD.SendPost(newIAM, voice)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, recoginzedText)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

}
