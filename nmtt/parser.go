package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

func (app *lineBotApp) parseApigwRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	requestBodyByte := []byte(r.Body)
	if !validateSignature(app.channelSecret, r.Headers["X-Line-Signature"], requestBodyByte) {
		return nil, linebot.ErrInvalidSignature
	}
	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal(requestBodyByte, request); err != nil {
		return nil, err
	}
	return request.Events, nil
}

// copied from github.com/line/line-bot-sdk-go/linebot/webhook.go
func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	hash.Write(body)
	return hmac.Equal(decoded, hash.Sum(nil))
}
