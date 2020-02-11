package yandexcld

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	folderID = "b1g22iiaq9tu0310a30e"
)

//ReadAudioFile - читаем файл
func ReadAudioFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("\n err %v \n", err)
	}

	return data
}

//RecognizeVoice - отправляем байты
func RecognizeVoice(iamToken string, voice []byte) (string, error) {
	endPoint := "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
	req, err := http.NewRequest("POST", endPoint+"?"+"topic=general&folderId="+folderID, bytes.NewBuffer(voice))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+iamToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

//SynthesizeVoice - синтезируем речь
func SynthesizeVoice(iamToken string, text string) ([]byte, error) {
	endPoint := "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"

	params := url.Values{
		"text":     {text},
		"lang":     {"ru-RU"},
		"folderId": {folderID},
	}

	req, err := http.NewRequest("POST", endPoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		fmt.Println("err", err)
	}

	req.Header.Set("Authorization", "Bearer "+iamToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
