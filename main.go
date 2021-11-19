package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Categorys struct {
	Userid   int
	Category string
}

var keyCategory = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		// tgbotapi.NewKeyboardButton("Выбрать категорию"),
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
		tgbotapi.NewKeyboardButton("7"),
		tgbotapi.NewKeyboardButton("8"),
		tgbotapi.NewKeyboardButton("9"),
	),
	tgbotapi.NewKeyboardButtonRow(
		// tgbotapi.NewKeyboardButton("Выбрать категорию"),
		tgbotapi.NewKeyboardButton("Создать категорию"),
		tgbotapi.NewKeyboardButton("Удалить категорию"),
	),
)

var keboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Все дела категории"),
		tgbotapi.NewKeyboardButton("Выбор категории"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Установить срок"),
		tgbotapi.NewKeyboardButton("Изменить статус"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Удалить дело"),
		tgbotapi.NewKeyboardButton("Удалить выполненные"),
	),
)

func main() {

	client, collectionTodos, colCategory := InitMongo()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	bot, updates := InitBot()

	allUserId := GetAllUserId(collectionTodos, "userid") //Получить userid всех пользователей со списками дел для уведомлений.

	//Сообщение об обновлении, отправляется один раз сразу после запуска сервера.
	//При перезапуске сервера удалить строчку НИЖЕ.
	SendUpdateNotification(allUserId, bot)

	//Отправка ежедневных уведоблений.
	SendNotification(allUserId, bot)

	flag := ""
	nameCategory := "Разное"

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//получение времени запроса
		now := FormatTime(update.Message.Time())
		userId := update.Message.From.ID
		getMessage := update.Message.Text
		msg := ""

		switch getMessage {
		case "/start":
			msg = fmt.Sprintf("Приветсвую тебя, <b><i>%s</i>! 🆕Релиз версии @mtdlBot 2.0:</b>\n- изменён интерфейс, добавлены клавиатуры для удобства использования\n- добавлены категории, теперь можно создавать/удалять свои категории дел\n- добавлена возможность установки срока выполнения дела\n- добавлен индикатор, отслеживающий истечение времени выполнения дела.\n\nЧтобы начать создавать свой список дел, просто напиши боту сообщение с новым делом. По умолчанию у всех пользователей создана только одна категория \"Разное\". Если у вас нет других категорий, все дела будут записываться в эту категорию. После создания новой категории, она сразу становится выбранной и новые дела будет сохраняться в неё. При удалении категории, все дела в ней будут так же удалены.", update.Message.From.FirstName)
			update := bson.M{"$set": bson.M{"category": "Разное"}}
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := colCategory.UpdateOne(ctx, Categorys{userId, "Разное"}, update, options.Update().SetUpsert(true))
			// _, err := colCategory.InsertOne(ctx, Categorys{userId, "Разное"})
			if err != nil {
				log.Fatal(err)
			}
		case "help":
			msg = "Тут будет содержаться справка."
		case "settings":
			msg = "Тут будут доступны настройки бота. Сейчас этот раздел в разработке."
		case "Удалить дело":
			msg = "Напишите номер удаляемого дела."
			flag = getMessage
		case "Изменить статус":
			msg = "Напишите номер дела."
			flag = getMessage
		case "Удалить выполненные":
			msg = CleanTodoList(collectionTodos, userId)
			msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
		case "Все дела категории":
			msg = PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
		case "Установить срок":
			msg = "Напишите номер дела и дату.\nФормат ввода: \"номер дела\". \"число\" \"месяц на англ. в сокращении\" \"год\" \"чч\":\"мм\"\nНапример --> <b>3. 12 Dec 2021 17:35</b>" //сделать ввод даты выполнения с кнопки?
			flag = getMessage
		case "Выбор категории":
			Msg := "Чтобы выбрать категорию, напишите её номер.\n"
			Msg += PrintCategory(GetAllUserCategory(colCategory, userId))
			m := tgbotapi.NewMessage(update.Message.Chat.ID, Msg)
			m.ParseMode = tgbotapi.ModeHTML
			m.ReplyMarkup = keyCategory
			bot.Send(m)
			flag = getMessage
		case "Создать категорию":
			msg = "Напишите название новой категории."
			flag = getMessage
		case "Удалить категорию":
			msg = "Напишите номер категории."
			flag = getMessage
		default:
			if flag == "Удалить дело" {
				msg = RemoveTodo(collectionTodos, userId, update.Message.Text, nameCategory)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
			} else if flag == "Изменить статус" {
				msg = ToggleTodo(collectionTodos, userId, update.Message.Text, nameCategory)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
			} else if flag == "Установить срок" {
				msg = Deadline(collectionTodos, userId, update.Message.Text, nameCategory)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
			} else if flag == "Выбор категории" {
				indexCategory := update.Message.Text
				result, i := ValidityOfIndex(colCategory, userId, indexCategory)
				if result {
					category := GetAllUserCategory(colCategory, userId)
					nameCategory = category[i-1].Category
					msg = fmt.Sprintf("Выбрана категория <b>%s</b>\n", nameCategory)
					msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
				} else {
					msg = "<i>❗Такая категория не существует.\n\n</i>"
				}
			} else if flag == "Создать категорию" {
				category := update.Message.Text
				if category != "Разное" && category != "разное" {
					ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
					_, err := colCategory.InsertOne(ctx, Categorys{userId, category})
					if err != nil {
						log.Fatal(err)
					}
					msg = "Категория создана."
				} else {
					msg = "<i>❗Такая категорию уже существует.</i>"
				}
			} else if flag == "Удалить категорию" {
				indexCategory := update.Message.Text
				msg = RemoveCategory(colCategory, collectionTodos, userId, indexCategory)
				nameCategory = "Разное"
			} else { //добавление нового дела
				msg = AddTodo(collectionTodos, userId, update.Message.Text, nameCategory, update.Message.Time()) //добавлять дела оп категориям
				msg += PrintTodoList(AllTodoList(collectionTodos, userId, nameCategory), now)
			}
			flag = ""
		}

		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		if flag == "" {
			msgParse.ReplyMarkup = keboard
		} else {
			msgParse.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
		bot.Send(msgParse)
	}
}
