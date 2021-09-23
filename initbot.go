package main

import (
	"io/ioutil"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetToken(filename string) string {
	token, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	tokenRep := strings.Replace(string(token), "\n", "", 1)
	return tokenRep
}

func GetURI(filename string) string {
	uri, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(uri)
}

func InitBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(GetToken("token.txt"))
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
