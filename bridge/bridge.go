package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Hello from bridge.go!")
	for {
		fmt.Println("Did this loop even run?")
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
			// proxy := taskData["proxy"].(string)
			// result = createTurboClient(proxy)
			fmt.Println("Creating a new client...")
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
