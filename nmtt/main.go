package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

type lineBotApp struct {
	channelSecret, accessToken string
	*linebot.Client
}

func createApp() *lineBotApp {
	channelSecret, accessToken := os.Getenv("CHANNEL_SECRET"), os.Getenv("ACCESS_TOKEN")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		log.Fatalln(err)
	}
	return &lineBotApp{channelSecret, accessToken, bot}
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	app := createApp()
	eventSlice, err := app.parseApigwRequest(request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Fatal("Invalid Signature")
		} else {
			log.Fatal(err)
		}
	}
	for _, event := range eventSlice {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var messageText string
				if strings.TrimSpace(message.Text) == "卓球したい" {
					messageText = "俺も卓球したいいい！！"
				}
				if messageText != "" {
					if _, err := app.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageText)).Do(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
