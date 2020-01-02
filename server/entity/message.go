package entity

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan SlackMessage)

type slackRtmStartResp struct {
	Ok  bool   `json:"ok"`
	URL string `json:"url"`
}

type SlackMessage struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

/*
MessageRoutine
Cretate goroutine which
  - Receive message from slack
  - Broadcast these messages to each clients
This method must be called at once
*/
func MessageRoutine(logger *log.Logger) {
	go readMessageFromslack(logger)
	go broadcastMessagesToClients(logger)
}

func ClientHandler(conn *websocket.Conn, logger *log.Logger) {
	clients[conn] = true

	for {
		var msg SlackMessage

		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			delete(clients, conn)
			break
		}

		logger.Println(msg)
	}
}

func readMessageFromslack(logger *log.Logger) {

	startURL := "https://slack.com/api/rtm.start"
	u, err := url.Parse(startURL)
	if err != nil {
		logger.Println(err)
		return
	}

	token := os.Getenv("SLACK_RTM_TOKEN")
	value := url.Values{}
	value.Add("token", token)

	u.RawQuery = value.Encode()
	logger.Println(u.String())

	url := u.String()

	resp, err := http.Get(url)
	if err != nil {
		logger.Println(err)
		return
	}
	defer resp.Body.Close()

	var d = &slackRtmStartResp{}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		logger.Println(err)
		return
	}

	c, _, err := websocket.DefaultDialer.Dial(d.URL, nil)
	if err != nil {
		logger.Fatal("dial:", err)
	}
	defer c.Close()

	for {
		_, messageByte, err := c.ReadMessage()
		if err != nil {
			logger.Println("read:", err)
			return
		}

		var message = SlackMessage{}
		if err := json.Unmarshal(messageByte, &message); err != nil {
			logger.Println("JSON Unmarshal error:", err)
			return
		}
		logger.Println(message)

		if message.Type == "message" {
			broadcast <- message
		}
	}
}

func broadcastMessagesToClients(logger *log.Logger) {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				logger.Printf("error occurred while writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
