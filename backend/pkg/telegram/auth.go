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

func ValidateInitDate(initDate string, botToken string) (int64, error)  {
	params, err := url.ParseQuery(initDate)
	if err != nil {
		return 0, err
	}

	hash := params.Get("hash")
	if hash == "" {
		return 0, errors.New("hash не найден в initData")
	}
	params.Del("hash")

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dateCheckStrings []string
	for _, k := range keys {
		dateCheckStrings = append(dateCheckStrings, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	dateCheckString := strings.Join(dateCheckStrings, "\n")

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secretKeyBytes := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secretKeyBytes)
	h.Write([]byte(dateCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))


	if calculatedHash != hash {
		return 0, errors.New("невалидная подпись initData (данные поддельные)")
	}

	var user WebAppUser
	userJSON := params.Get("user")
	if userJSON == "" {
        return 0, errors.New("user не найден в initData")
    }

    err = json.Unmarshal([]byte(userJSON), &user)
    if err != nil {
        return 0, err
    }

    return user.ID, nil

}