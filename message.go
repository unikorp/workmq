package workmq

import (
	"encoding/json"
	"fmt"
)

// Message struct
type Message struct {
	Queue string `json:"queue"`
	Body  string `json:"body"`
}

// TransformStringToMessage transforms a string value to a Message struct
func TransformStringToMessage(value []byte) Message {
	message := Message{}
	err := json.Unmarshal(value, &message)

	if err != nil {
		fmt.Println("Unable to transform string to Message struct:", err)
	}

	return message
}
