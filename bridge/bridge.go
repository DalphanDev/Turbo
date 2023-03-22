package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DalphanDev/Turbo/src"
	"github.com/google/uuid"
)

type ClientResponse struct {
	Command  string `json:"command"`
	ClientID string `json:"clientID"`
}

type DoResponse struct {
	Command  string `json:"command"`
	ClientID string `json:"clientID"`
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
			proxy := taskData["proxy"].(string)
			client := src.NewTurboClient(proxy)
			clientID := uuid.New().String()
			clients[clientID] = client
			result = ClientResponse{
				Command:  command,
				ClientID: clientID,
			}
			// fmt.Println(result)

		case "do":
			// fmt.Println("Sending a request...")
			// url := taskData["url"].(string)
			// method := taskData["method"].(string)
			// headers := taskData["headers"].(string)
			// body := taskData["body"].(string)
			clientID := taskData["clientID"].(string)
			// fmt.Println(url)
			// fmt.Println(method)
			// fmt.Println(headers)
			// fmt.Println(body)
			// fmt.Println(clientID)

			result = DoResponse{
				Command:  command,
				ClientID: clientID,
			}
			// fmt.Println(result)

			// options := src.RequestOptions{
			// 	URL:     "https://eoobxe7m89qj9cl.m.pipedream.net",
			// 	Headers: nil,
			// 	Body:    strings.NewReader(body), // Can either use nil or a string reader.
			// }

			// resp, err := client.Do(method, options)

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
