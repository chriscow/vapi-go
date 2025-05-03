package vapi

/*
curl -X POST 'https://aws-us-west-2-production1-phone-call-websocket.vapi.ai/7420f27a-30fd-4f49-a995-5549ae7cc00d/control'
-H 'content-type: application/json'
--data-raw '{
  "type": "add-message",
  "message": {
    "role": "system",
    "content": "New message added to conversation"
  },
  "triggerResponseEnabled": true
}'
*/

type CallControlAddMessage struct {
	Type                   string        `json:"type"`
	TriggerResponseEnabled bool          `json:"triggerResponseEnabled"`
	Message                OpenAIMessage `json:"message"`
}
