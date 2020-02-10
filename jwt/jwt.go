package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	keyID            = "ajelgl72fg5hpoi9m7ue"
	serviceAccountID = "ajeuuablhh36r6a9fmrg"
	keyFile          = "jwt/private.pem"
)

// Формирование JWT.
func signedToken() string {
	issuedAt := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.StandardClaims{
		Issuer:    serviceAccountID,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: issuedAt.Add(time.Hour).Unix(),
		Audience:  "https://iam.api.cloud.yandex.net/iam/v1/tokens",
	})
	token.Header["kid"] = keyID

	privateKey := loadPrivateKey()
	signed, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signed
}

func loadPrivateKey() *rsa.PrivateKey {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
	return rsaPrivateKey
}

//GetNewIAMToken - получение нового IAM по вновь сформированному JWT
func GetNewIAMToken() string {
	jot := signedToken()
	fmt.Println("My JWT", jot)
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("%s: %s", resp.Status, body))
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data.IAMToken
}
