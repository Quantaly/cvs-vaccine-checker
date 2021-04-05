package structs

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Application

type Persistence struct {
	LastTimestamp string
	HadVaccines   bool
}

// CVS

type Response struct {
	ResponsePayloadData Payload
	ResponseMetaData    Metadata
}

type Payload struct {
	CurrentTime        string
	Data               Data
	IsBookingCompleted bool
}

type Data struct {
	CO []City
}

type City struct {
	City   string
	State  string
	Status string
}

type Metadata struct {
	StatusDesc     string
	ConversationId string
	RefId          string
	Operation      string
	Version        string
	StatusCode     string
}

// Discord

type Message struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

func (m Message) Send() error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = http.Post(discordWebhookUrl, "application/json", bytes.NewReader(body))
	return err
}

type Embed struct {
	Title       string  `json:"title"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	Timestamp   string  `json:"timestamp"`
	Color       uint32  `json:"color"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
