package main

import (
	"bytes"
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
	PAGE_TOKEN    = os.Getenv("PAGE_TOKEN")
	VERIFY_TOKEN  = "developers-are-gods"
	FB_APP_SECRET = os.Getenv("FB_APP_SECRET")
	AUTH_TOKEN    = os.Getenv("AUTH_TOKEN")
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

type VisualSearchResponse struct {
	Links []string `json:"links"`
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

func visual_search(queryImage string) *VisualSearchResponse {
	url := "http://localhost:8000/quick"
	log.Info(queryImage)

	q := map[string]string{"image_url": queryImage}

	jsonStr, err := json.Marshal(q)
	if err != nil {
		fmt.Println("error:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	log.Info(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Info("Visual search server response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	links := new(VisualSearchResponse)
	err = json.Unmarshal(body, &links)

	return links
}

func main() {
	bot := mbotapi.NewBotAPI(PAGE_TOKEN, VERIFY_TOKEN, FB_APP_SECRET)

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

		// Send messages or send image results
		if len(callback.Message.Attachments) == 0 {
			bot.Send(callback.Sender, msg, mbotapi.RegularNotif)
		} else {
			log.Info(callback.Message.Attachments[0].Payload.URL)
			products, err := callShopStyle()
			if err != nil {
				log.Error(err)
			}

			template := mbotapi.NewGenericTemplate()

			i := 0
			for _, product := range products.Products {
				element := mbotapi.NewElement(product.Name)
				template.Elements = append(template.Elements, element)
				template.Elements[i].ImageURL = product.Image.Sizes.Best.URL
				buyButton := mbotapi.NewURLButton("Buy it!", product.ClickURL)
				template.Elements[i].Buttons = append(template.Elements[i].Buttons, buyButton)
				i = i + 1
			}
			bot.Send(callback.Sender, template, mbotapi.RegularNotif)
		}
	}
}
