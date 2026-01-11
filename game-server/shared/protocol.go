package shared

import (
	"encoding/json"
	"fmt"
)

type Action struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

const (
	ActionTypeMove = "move"
	ActionTypeChat = "chat"
)

type MoveData struct {
	Direction string `json:"direction"`
}

type ChatData struct {
	Message string `json:"message"`
}

func ParseAction(data []byte) (*Action, error) {
	var action Action
	if err := json.Unmarshal(data, &action); err != nil {
		return nil, fmt.Errorf("Invalid action: %w", err)
	}

	if action.Type == "" {
		return nil, fmt.Errorf("Invalid action type")
	}

	return &action, nil
}

func (a *Action) GetMoveData() (*MoveData, error) {
	if a.Type != ActionTypeMove {
		return nil, fmt.Errorf("Invalid action type")
	}

	var moveData MoveData
	if err := json.Unmarshal(a.Data, &moveData); err != nil {
		return nil, fmt.Errorf("Invalid move data: %w", err)
	}

	return &moveData, nil
}

func (a *Action) GetChatData() (*ChatData, error) {
	if a.Type != ActionTypeChat {
		return nil, fmt.Errorf("Invalid action type")
	}

	var chatData ChatData
	if err := json.Unmarshal(a.Data, &chatData); err != nil {
		return nil, fmt.Errorf("Invalid chat data: %w", err)
	}

	return &chatData, nil
}
