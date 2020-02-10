package main

import (
	jwt "github.com/vizzyolog/showcase/jwt"
	yandexCLD "github.com/vizzyolog/showcase/yandexCLD"
)

const (
	botURL = "t.me/Anti_twaddle_bot"
)

func openConfig() {

}

func main() {
	newIAM := jwt.GetNewIAMToken()

	voice := yandexCLD.ReadAudioFile("sources/speech3.ogg")
	yandexCLD.SendPost(newIAM, voice)

}
