package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log" // <-- Добавь этот импорт
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

func ValidateInitData(initData string, botToken string) (int64, string, error) {
	// 🔥 ОТЛАДКА: Смотрим, что именно прислал Telegram (первые 150 символов)
	log.Printf("🔍 RAW initData: %s", initData[:min(len(initData), 150)])

	params, err := url.ParseQuery(initData)
	if err != nil {
		return 0, "", err
	}

	// 🔥 ОТЛАДКА: Печатаем все ключи, которые пришли от Telegram
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	log.Printf("🔍 Все ключи в initData: %v", keys)

	hash := params.Get("hash")
	if hash == "" {
		return 0, "", errors.New("hash не найден в initData")
	}
	params.Del("hash")

	sort.Strings(keys)
	// Убираем 'hash' из ключей для проверки, так как мы его уже удалили из params, 
	// но keys мы собрали ДО удаления. Пересоберем keys без hash:
	var dataCheckKeys []string
	for k := range params {
		dataCheckKeys = append(dataCheckKeys, k)
	}
	sort.Strings(dataCheckKeys)

	var dataCheckStrings []string
	for _, k := range dataCheckKeys {
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

	// 🔥 ДОСТАЕМ start_param
	startParam := params.Get("start_param")
	log.Printf("🔍 Извлеченный start_param: '%s'", startParam) // <-- Пустые кавычки означают, что Telegram его не прислал

	return user.ID, startParam, nil
}

// Вспомогательная функция для min (если используешь Go < 1.21)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}