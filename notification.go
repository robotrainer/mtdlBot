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
		msg := "<b>🆕Релиз версии @mtdlBot 2.0:\n- изменён интерфейс, добавлены клавиатуры для удобства использования\n- добавлены категории, теперь можно создавать/удалять свои категории дел\n- добавлена возможность установки срока выполнения дела\n- добавлен индикатор, отслеживающий истечение времени выполнения дела.\n\nЧтобы начать создавать свой список дел, просто напиши боту сообщение с новым делом. По умолчанию у всех пользователей создана только одна категория \"Разное\". Если у вас нет других категорий, все дела будут записываться в эту категорию. После создания новой категории, она сразу становится выбранной и новые дела будет сохраняться в неё. При удалении категории, все дела в ней будут так же удалены.</b>"
		msgParse := tgbotapi.NewMessage(chatId, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
