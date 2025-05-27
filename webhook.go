package vapi

import (
	"time"
)

// Message types
const (
	MsgTypeAssistantRequest   = "assistant-request"
	MsgTypeToolCalls          = "tool-calls"
	MsgTypeTransferDestReq    = "transfer-destination-request"
	MsgTypeConversationUpdate = "conversation-update"
	MsgTypeEndOfCallReport    = "end-of-call-report"
	MsgTypeFunctionCall       = "function-call"
	MsgTypeHang               = "hang"
	MsgTypeSpeechUpdate       = "speech-update"
	MsgTypeStatusUpdate       = "status-update"
	MsgTypeUserInterrupted    = "user-interrupted"
)

// Message represents the incoming message from VAPI webhooks
type WebhookMessage struct {
	Message struct {
		Type string `json:"type"`
	} `json:"message"`
}

type AssistantRequest struct {
	Type        string       `json:"type"`
	PhoneNumber *PhoneNumber `json:"phoneNumber,omitempty"`
	Timestamp   *float64     `json:"timestamp,omitempty"`
	Artifact    *Artifact    `json:"artifact,omitempty"`
	Assistant   *Assistant   `json:"assistant,omitempty"`
	Customer    *Customer    `json:"customer,omitempty"`
	Call        *Call        `json:"call,omitempty"`
}

type AssistantRequestEnvelope struct {
	AssistantRequest AssistantRequest `json:"message"`
}

type AssistantRequestResponse struct {
	Destination        *Destination `json:"destination,omitempty"`
	AssistantId        *string      `json:"assistantId,omitempty"`
	Assistant          *Assistant   `json:"assistant,omitempty"`
	AssistantOverrides *Assistant   `json:"assistantOverrides,omitempty"`
	CustomerID         *string      `json:"customerId,omitempty"`
	Customer           *Customer    `json:"customer,omitempty"`
}

type EndOfCallReport struct {
	ID        *string   `json:"id"`
	Timestamp *float64  `json:"timestamp,omitempty"`
	Type      string    `json:"type"`
	Artifact  *Artifact `json:"artifact"`
	Analysis  *Analysis `json:"analysis"`

	// Optional fields
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	EndedAt     *time.Time `json:"endedAt,omitempty"`
	EndedReason string     `json:"endedReason"`
	Cost        *float64   `json:"cost,omitempty"`
	Costs       []Cost     `json:"costs,omitempty"`

	Summary    *string `json:"summary,omitempty"`
	Transcript *string `json:"transcript,omitempty"`

	Messages           []Message    `json:"messages,omitempty"`
	RecordingUrl       *string      `json:"recordingUrl,omitempty"`
	StereoRecordingUrl *string      `json:"stereoRecordingUrl,omitempty"`
	Call               *Call        `json:"call,omitempty"`
	PhoneNumber        *PhoneNumber `json:"phoneNumber,omitempty"`
	Customer           *Customer    `json:"customer,omitempty"`

	Assistant *Assistant `json:"assistant,omitempty"`
}

// EndOfCallReportEnvelope represents the report generated at the end of a call
type EndOfCallReportEnvelope struct {
	EndOfCallReport EndOfCallReport `json:"message"`
}

type ConversationUpdateEnvelope struct {
	ConversationUpdate ConversationUpdate `json:"message"`
}

type ConversationUpdate struct {
	Type           string          `json:"type"`
	OpenAIMessages []OpenAIMessage `json:"messagesOpenAIFormatted"`
	Messages       []Message       `json:"messages,omitempty"`
	PhoneNumber    *PhoneNumber    `json:"phoneNumber,omitempty"`
	CustomerID     *string         `json:"customerId,omitempty"`
	Customer       *Customer       `json:"customer,omitempty"`
	Call           *Call           `json:"call,omitempty"`
}
