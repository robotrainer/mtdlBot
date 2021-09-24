package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetURI(filename string) string {
	uri, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	uriRep := strings.Replace(string(uri), "\n", "", 1)
	return uriRep
}

func InitMongo() (*mongo.Client, *mongo.Collection) {
	//Создаём нового клиента базы данных с поключением по указанному URL
	client, err := mongo.NewClient(options.Client().
		ApplyURI(GetURI("uri.txt")))
	if err != nil {
		log.Fatal(err)
	}

	//Подключаемся к базе данных
	fmt.Println("Client connecting...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //ждём ответа в течении 10 сек
	err = client.Connect(ctx)                                           //если ответа нет, вернёт ошибку
	if err != nil {
		log.Fatal(err)
	}

	// defer client.Disconnect(ctx) //откладываем момент отключения от базы данных
	//проверяем подклчение
	fmt.Println("Ping connecting...")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second) //ждём ответа в течении 10 сек
	err = client.Ping(ctx, readpref.Primary())                         //отправляет сигнал ping, чтобы проверить, может ли клиент быть подключен к базе данных
	if err != nil {
		log.Fatal(err)
	}

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(databases)

	collection := client.Database("mtdlBot").Collection("todolist")

	return client, collection
}
