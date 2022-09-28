package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := getConfigurations()
	if err != nil {
		fmt.Println(err)
		return
	}
	lastNotificationTime := time.Time{}
	tickId := 1
	for range time.Tick(time.Second * time.Duration(config.App.RefreshFrequencySeconds)) {
		outputStr := "tick " + strconv.Itoa(tickId) + "| "
		tickId += 1
		outputStr += "time: " + time.Now().Format("2006-01-02 15:04:05.999999999") + "| "
		postRes := getPrices(config.Endpoints, config.App.RequestTimeoutDuration)

		hasErr := false
		for _, res := range postRes {
			if res.Err != nil {
				outputStr += "Endpoint " + res.Endpoint.Name + " error. " + res.Err.Error() + "| "
				hasErr = true
			} else {
				outputStr += res.Endpoint.Name + ": " + fmt.Sprintf("%f", res.Endpoint.ParsedPrice) + "| "
			}
		}
		if hasErr {
			fmt.Println(outputStr)
			continue
		}
		diffPercentage := diffInPercentage(postRes[0].Endpoint.ParsedPrice, postRes[1].Endpoint.ParsedPrice)
		outputStr += "diffPercentage: " + fmt.Sprintf("%.2f", diffPercentage)
		fmt.Println(outputStr)
		if time.Since(lastNotificationTime).Seconds() >= config.App.NotificationTimeoutSeconds &&
			diffPercentage != 0 &&
			diffPercentage > config.App.NotifyPercentage {
			sendMsgTelegram(outputStr, config.App.TelegramBotAPIKey, config.App.TelegramChatId)
			lastNotificationTime = time.Now()
		}
	}
}

func sendMsgTelegram(messageText string, tgBotAPIKey string, chatID int64) error {
	bot, err := tgbotapi.NewBotAPI(tgBotAPIKey)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, messageText)
	_, err = bot.Send(msg)
	return err
}

func diffInPercentage(firstNumber float64, secondNumber float64) float64 {
	if firstNumber == 0 || secondNumber == 0 {
		return 0
	}
	if firstNumber > secondNumber {
		return ((firstNumber - secondNumber) / math.Abs(secondNumber)) * 100
	} else {
		return ((secondNumber - firstNumber) / math.Abs(firstNumber)) * 100
	}
}
