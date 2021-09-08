package main

type Update struct {
	Message  Message `json:"message"`
	UpdateId int     `json:"update_id"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	From User   `json:"from"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type User struct {
	Id int `json:"id"`
}
