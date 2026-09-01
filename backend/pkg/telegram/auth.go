package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type WebAppUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// 🔥 ВОЗВРАЩАЕМ ВСЮ СТРУКТУРУ ПОЛЬЗОВАТЕЛЯ
func ValidateInitData(initData string, botToken string) (WebAppUser, string, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return WebAppUser{}, "", err
	}

	hash := params.Get("hash")
	if hash == "" {
		return WebAppUser{}, "", errors.New("hash не найден")
	}
	params.Del("hash")

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckStrings []string
	for _, k := range keys {
		dataCheckStrings = append(dataCheckStrings, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secretKeyBytes := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secretKeyBytes)
	h.Write([]byte(strings.Join(dataCheckStrings, "\n")))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != hash {
		return WebAppUser{}, "", errors.New("невалидная подпись")
	}

	var user WebAppUser
	userJSON := params.Get("user")
	if userJSON == "" {
		return WebAppUser{}, "", errors.New("user не найден")
	}

	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return WebAppUser{}, "", err
	}

	startParam := params.Get("start_param")
	return user, startParam, nil
}