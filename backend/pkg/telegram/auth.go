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

// ValidateInitData возвращает (tgID, startParam, error)
func ValidateInitData(initData string, botToken string) (int64, string, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return 0, "", err
	}

	hash := params.Get("hash")
	if hash == "" {
		return 0, "", errors.New("hash не найден в initData")
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
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secretKeyBytes := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secretKeyBytes)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != hash {
		return 0, "", errors.New("невалидная подпись initData")
	}

	var user WebAppUser
	userJSON := params.Get("user")
	if userJSON == "" {
		return 0, "", errors.New("user не найден в initData")
	}

	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return 0, "", err
	}

	// 🔥 ДОСТАЕМ start_param (invite_link для клиента)
	startParam := params.Get("start_param")

	return user.ID, startParam, nil
}