package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUserId(collection *mongo.Collection, fieldName string) []interface{} {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	allUserId, err := collection.Distinct(ctx, fieldName, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(allUserId...)
	return allUserId
}

// Вызов переданной функции в указанное время.
func CallAtTime(hour, min, sec int, bot *tgbotapi.BotAPI, chatId int64, f func(bot *tgbotapi.BotAPI, chatId int64)) error {
	//ПОлучаем локацию Москва
	msk, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	// Вычисляем время первого запуска.
	now := time.Now().Local()
	t := now.In(msk)
	firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), hour, min, sec, 0, msk)
	if firstCallTime.Before(now) {
		// Если получилось время раньше текущего, прибавляем сутки.
		firstCallTime = firstCallTime.Add(time.Hour * 24)
	}

	// Вычисляем временной промежуток до запуска.
	duration := firstCallTime.Sub(time.Now().Local())

	fmt.Printf("Время сервера: %v\nВремя первого запуска по МСК: %v\nВремя до первого уведомления: %v\n Текущее время в МСК: %v\n", now, firstCallTime, duration, t)

	//Горутина вызова функции через указанный период.
	go func() {
		time.Sleep(duration)
		for {
			f(bot, chatId)
			// Следующий запуск через сутки.
			time.Sleep(time.Hour * 24)
		}
	}()

	return nil
}

// Функция уведомления.
func Notification(bot *tgbotapi.BotAPI, chatId int64) {
	msg := "<i><b>🔔Уведомление</b></i>\nНе забудь выполнить запланированные дела и записать новые, чтобы всегда держать список под рукой!"
	msgParse := tgbotapi.NewMessage(chatId, msg)
	msgParse.ParseMode = tgbotapi.ModeHTML
	bot.Send(msgParse)
	fmt.Printf("Сообщение успешно отправлено!\nВремя отправления: %v\n", time.Now())
}

//Функция отправки уведомлений.
func SendNotification(allUserId []interface{}, bot *tgbotapi.BotAPI) {
	for _, id := range allUserId {
		chatId := int64(id.(int32))
		err := CallAtTime(0, 0, 0, bot, chatId, Notification) //сделать рассылку сообщиния по Chat.ID из базы данных
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

func SendUpdateNotification(allUserId []interface{}, bot *tgbotapi.BotAPI) {
	for _, id := range allUserId {
		chatId := int64(id.(int32))
		msg := "<b>🆕Обновление @mtdlBot №1:</b>\n1. Для добавления нового дела теперь не надо вводить команду /add, просто напиши дело и отправь его боту.\n2. При вводе команды /start, бот приветсвует пользователя.\n3. Добавлена команда /help со списком доступных возможностей бота.\n4. Добавлены уведомления об обновлении.\n5. Добавлены ежедневные уведомления с напоминанием о выполнении и сохранении дел, время рассылки 00:00 по МСК.\n\n⚜️<b>Анонс плана обновлений:</b>\n- категории дел;\n- новые статусы выполнения дел (в процессе, неактуально и т.п.);\n- установка срока выполнения дела."
		msgParse := tgbotapi.NewMessage(chatId, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
