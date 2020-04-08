package main

import "encoding/json"

type message struct {
	UniqueId     string `json:"uniqueId"`
	Text         string `json:"text"`
	Action       string `json:"action"`
	UserId       string `json:"userId"`
	UserUsername string `json:"userUsername"`
	Room         string `json:"room"`
}

// toJson converts message to json object, hides the sensitive info if hideSensitive is true
func (message *message) toJson(hideSensitive bool) ([]byte, error) {
	var (
		jsonData  []byte
		jsonError error
	)

	if hideSensitive {
		safeData := map[string]interface{}{"action": message.Action,
			"username": message.UserUsername, "text": message.Text}
		jsonData, jsonError = json.Marshal(safeData)
	} else {
		jsonData, jsonError = json.Marshal(message)
	}

	if jsonError != nil {
		return nil, jsonError
	}
	return jsonData, nil
}

func (message *message) toProtectedMap() map[string]interface{}{
	return map[string]interface{}{"action": message.Action,
		"username": message.UserUsername, "text": message.Text, "id": message.UniqueId}
}

func (message *message) setId() {
	message.UniqueId = randId()
}

func messageFromJson(jsonPayload []byte) (error, *message) {
	var model message
	err := json.Unmarshal(jsonPayload, &model)
	if err != nil {
		return err, nil
	}
	model.setId()
	return nil, &model
}

