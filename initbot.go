package main

import (
	"io/ioutil"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const token = "token.txt"

func GetToken(filename string) string {
	token, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	tokenRep := strings.Replace(string(token), "\n", "", 1)
	return tokenRep
}

func InitBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(GetToken(token))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	return bot, updates
}
