package main

import (
	"log"
	"os"
	"strings"

	"github.com/yo-tak/nmtt/gym"

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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	app := createApp()
	eventSlice, err := app.parseApigwRequest(request)
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range eventSlice {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageWithoutSpace := strings.TrimSpace(message.Text)
				var messageText string
				if strings.Contains(messageWithoutSpace, "卓球したい") {
					messageText = "俺も卓球したいいい！！"
				} else if strings.Contains(messageWithoutSpace, "卓球したくない") {
					messageText = "俺は卓球したい！！！"
				} else if messageWithoutSpace == "中野の予定" {
					log.Println("let's scrape schedule")
					url, err := gym.GetNakanoCurrentNotification()
					if err != nil {
						log.Fatal(err)
					}
					log.Println("it seems there were no error here")
					messageText = "直近の予定表だよ：\n" + url
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
	lambda.Start(handler)
}
