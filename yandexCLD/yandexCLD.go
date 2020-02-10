package yandexcld

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	folderID = "b1g22iiaq9tu0310a30e"
)

func ReadAudioFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("\n err %v \n", err)
	}

	return data
}

//SendPost - отправляем байты
func SendPost(iamToken string, voice []byte) string {
	url := "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
	req, err := http.NewRequest("POST", url+"?"+"topic=general&folderId="+folderID, bytes.NewBuffer(voice))
	if err != nil {
		fmt.Println("err", err)
	}

	req.Header.Set("Authorization", "Bearer "+iamToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(bodyBytes)

}
