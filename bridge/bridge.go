package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/DalphanDev/Turbo/http"
	"github.com/DalphanDev/Turbo/src"
	"github.com/google/uuid"
)

type ClientResponse struct {
	Command  string `json:"command"`
	ClientID string `json:"clientID"`
}

type DoResponse struct {
	Command    string      `json:"command"`
	ClientID   string      `json:"clientID"`
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
}

func main() {
	inputReader := bufio.NewReader(os.Stdin)

	clients := make(map[string]*src.TurboClient)
	for {
		inputData, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v", err)
		}

		var taskData map[string]interface{}
		err = json.Unmarshal([]byte(inputData), &taskData)
		if err != nil {
			fmt.Printf("Error unmarshaling task data: %v", err)
		}

		command := taskData["command"].(string)

		var result interface{}

		switch command {
		case "new_client":
			// Creating an turbo client...
			proxy := taskData["proxy"].(string)
			client := src.NewTurboClient(proxy)
			clientID := uuid.New().String()
			clients[clientID] = client
			result = ClientResponse{
				Command:  command,
				ClientID: clientID,
			}

		case "do":
			// Sending a request...
			url := taskData["url"].(string)
			method := taskData["method"].(string)
			// headers := taskData["headers"].(string)
			body := taskData["body"].(string)
			clientID := taskData["clientID"].(string)

			myClient := clients[clientID]

			options := src.RequestOptions{
				URL:     url,
				Headers: nil,
				Body:    strings.NewReader(body), // Can either use nil or a string reader.
			}

			resp, err := myClient.Do(method, options)
			if err != nil {
				panic(err)
			}

			result = DoResponse{
				Command:    command,
				ClientID:   clientID,
				StatusCode: resp.StatusCode,
				Headers:    resp.Headers,
				Body:       resp.Body,
			}

		default:
			fmt.Printf("Invalid command: %s", command)
		}

		outputData, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Error marshaling response: %v", err)
		}
		fmt.Println(string(outputData))
	}
}
