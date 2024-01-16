package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var app *firebase.App = nil

func Init(ctx context.Context) (*firebase.App, error) {
	conf := &firebase.Config{
		DatabaseURL: "https://lasm-go-default-rtdb.europe-west1.firebasedatabase.app/",
	}
	opt := option.WithCredentialsFile("credentials/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	fmt.Println("Firebase initialized")
	return app, err
}

type Log struct {
	Endpoint  string `json:"endpoint"`
	Ip        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Client    string `json:"client"`
	Status    string `json:"status"`
}

func SendLog(l ...Log) error {
	ctx := context.Background()
	if app == nil {
		var err error
		app, err = Init(ctx)
		if err != nil {
			return err
		}
	}
	client, err := app.Database(ctx)
	if err != nil {
		return err
	}

	ref := client.NewRef("/log")
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
		return err
	}
	// if l is nil make a default log
	if len(l) == 0 {
		l = append(l,
			Log{
				Endpoint:  "default",
				Ip:        "127.0.0.1",
				Timestamp: "2021-05-01 00:00:00",
				Client:    "default",
				Status:    "success",
			})
	}
	for i := range l {
		_, err = ref.Push(ctx, l[i])
	}
	if err != nil {
		log.Fatalln("Error setting value:", err)
		return err
	}
	return nil
}

func SendFCM(token string, title string, body string) error {
	ctx := context.Background()
	if app == nil {
		var err error
		app, err = Init(ctx)
		if err != nil {
			return err
		}
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return err
	}

	message := &messaging.Message{
		Data: map[string]string{
			"title": title,
			"body":  body,
		},
		Token: token,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalf("error sending message: %v - %s\n", err, response)
		return err
	}

	return nil
}
