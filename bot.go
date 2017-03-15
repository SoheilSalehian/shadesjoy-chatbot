package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/abhinavdahiya/go-messenger-bot"
)

var (
	PAGE_TOKEN   = os.Getenv("PAGE_TOKEN")
	VERIFY_TOKEN = "developers-are-gods"
	AUTH_TOKEN   = os.Getenv("AUTH_TOKEN")
)

type ApiAiInput struct {
	Status struct {
		Code      int
		ErrorType string
	}
	Result struct {
		Action           *string
		ActionIncomplete bool
		Speech           string
	} `json:"result"`
}

func main() {
	bot := mbotapi.NewBotAPI(PAGE_TOKEN, VERIFY_TOKEN, "a5318d65de57385a49bd90bfc9efc31c")

	callbacks, mux := bot.SetWebhook("/webhook")
	go http.ListenAndServe("0.0.0.0:9091", mux)
	log.Info("starting server on :9091")

	var msg interface{}
	for callback := range callbacks {
		log.Printf("[%#v] %s", callback.Sender, callback.Message.Text)

		if resp, err := getApiAiResponse(callback.Message.Text, callback.Sender.ID); err == nil {
			msg = mbotapi.NewMessage(resp)
		} else {
			msg = mbotapi.NewMessage(callback.Message.Text)
		}
		bot.Send(callback.Sender, msg, mbotapi.RegularNotif)
	}
}

func getApiAiResponse(message string, senderId int64) (resp string, err error) {
	params := url.Values{}
	params.Add("query", message)
	params.Set("sessionId", string(senderId))

	url := fmt.Sprintf("https://api.api.ai/v1/query?V=20160518&lang=En&%s", params.Encode())
	ai, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	ai.Header.Set("Authorization", "Bearer "+AUTH_TOKEN)

	if resp, err := http.DefaultClient.Do(ai); err != nil {
		return "", err
	} else {
		defer resp.Body.Close()

		var input ApiAiInput
		datastring, _ := ioutil.ReadAll(resp.Body)
		err := json.NewDecoder(strings.NewReader(string(datastring))).Decode(&input)
		if err != nil {
			return "", err
		}

		return input.Result.Speech, nil
	}
}
