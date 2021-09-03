package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		word := strings.Fields(update.Message.Text)
		fmt.Println(word)

		switch word[0] {
		case "ADD":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда ADD"))
		case "RM":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда RM"))
		case "TG":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда TG"))
		case "CL":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда CL"))
		case "ALL":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда ALL"))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
		}
	}
}
