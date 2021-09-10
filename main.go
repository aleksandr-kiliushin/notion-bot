package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// app entry point
func main() {
	botUrl := "https://api.telegram.org/bot" + getEnvVariable("TELEGRAM_NOTION_BOT_TOKEN")

	offset := 0

	for {
		updates, err := getUpdates(botUrl, offset)

		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}

		for _, update := range updates {
			addTodo(botUrl, update)
			offset = update.UpdateId + 1
		}

		// fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var restResponse RestResponse

	err = json.Unmarshal(body, &restResponse)

	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func addTodo(botUrl string, update Update) error {
	telegramUserId, err := strconv.Atoi(getEnvVariable("TELEGRAM_USER_ID"))

	if err != nil {
		return nil
	}

	if update.Message.From.Id != telegramUserId {
		return nil
	}

	var jsonStr = []byte(fmt.Sprintf(`{
		"parent": {
			"database_id": "%s"
		},
		"properties": {
			"title": {
				"title": [
					{
						"text": {
							"content": "%s"
						}
					}
				]
			},
			"Status":{
				"id":"F%3A%5CQ",
				"type":"select",
				"select":{
					"id":"1",
					"name":"Not started",
					"color":"red"
				}
			}
		}
	}`, getEnvVariable("BOARD_DATABASE_ID"), update.Message.Text))

	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "Bearer "+getEnvVariable("INTERNAL_INTEGRATION_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2021-08-16")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	return nil
}
