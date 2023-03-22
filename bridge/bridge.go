package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DalphanDev/Turbo/src"
)

type ClientResponse struct {
	Command string `json:"command"`
	Client  string `json:"client"`
}

func main() {
	inputReader := bufio.NewReader(os.Stdin)
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

		fmt.Println(command)
		var result interface{}

		switch command {
		case "new_client":
			fmt.Println("Creating a new client...")
			proxy := taskData["proxy"].(string)
			client := src.NewTurboClient(proxy)
			stringifiedClient := fmt.Sprintf("%p", client)
			result = ClientResponse{
				Command: command,
				Client:  stringifiedClient,
			}
			fmt.Println(result)

		// Add more cases for other functions as needed
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
