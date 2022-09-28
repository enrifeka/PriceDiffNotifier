package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

func getConfigurations() (*Configuration, error) {
	c := Configuration{}
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		return nil, err
	}
	if c.App.NotifyPercentage < 0 || c.App.NotifyPercentage > 100 {
		return nil, errors.New("notifyOnPercentageDifference must be from 0 to 100")
	}
	if c.App.RefreshFrequencySeconds <= 0 {
		return nil, errors.New("refreshFrequencySeconds must be greater than 0")
	}
	if c.App.RequestTimeoutSeconds <= 0 {
		return nil, errors.New("requestTimeoutSeconds must be greater than 0")
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%fs", c.App.RequestTimeoutSeconds))
	if err != nil {
		return nil, errors.New("cannot convert requestTimeoutSeconds to duration " + err.Error())
	}
	c.App.RequestTimeoutDuration = duration
	if c.App.NotificationTimeoutSeconds <= 0 {
		return nil, errors.New("notificationTimeoutSeconds must be greater than 0")
	}
	// get only active endpoints
	var activeEndpoints []Endpoint
	for _, e := range c.Endpoints {
		if e.Active {
			activeEndpoints = append(activeEndpoints, e)
		}
	}
	// 2 active endpoints are required
	if len(activeEndpoints) != 2 {
		return nil, errors.New("two active endpoints are required")
	}
	c.Endpoints = activeEndpoints
	return &c, nil
}

type Configuration struct {
	App       AppConfig  `json:"appConfig"`
	Endpoints []Endpoint `json:"endpoints"`
}

type AppConfig struct {
	TelegramBotAPIKey          string  `json:"telegramBotAPIKey"`
	TelegramChatId             int64   `json:"telegramChatId"`
	NotifyPercentage           float64 `json:"notifyOnPercentageDifference"`
	RefreshFrequencySeconds    float64 `json:"refreshFrequencySeconds"`
	RequestTimeoutSeconds      float64 `json:"requestTimeoutSeconds"`
	RequestTimeoutDuration     time.Duration
	NotificationTimeoutSeconds float64 `json:"notificationTimeoutSeconds"`
}

type Endpoint struct {
	Name         string `json:"name"`
	PriceAPI     string `json:"priceAPI"`
	PriceTagPath string `json:"priceTagPath"`
	Active       bool   `json:"active"`
	ParsedPrice  float64
}
