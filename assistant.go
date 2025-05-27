package vapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chriscow/minds"
)

// Alex: 0f7534fa-deee-4feb-a48f-6e6e64eb38e7

// Assistant represents the top-level configuration for a voice interaction
type Assistant struct {
	ID                           *string                `json:"id,omitempty"`
	Name                         *string                `json:"name,omitempty"`
	Voice                        *ElevenLabsVoiceConfig `json:"voice,omitempty"`
	Model                        *ModelConfig           `json:"model,omitempty"`
	Transcriber                  *TranscriberConfig     `json:"transcriber,omitempty"`
	FirstMessage                 *string                `json:"firstMessage,omitempty"`
	ClientMessages               []string               `json:"clientMessages,omitempty"`
	ServerMessages               []string               `json:"serverMessages,omitempty"`
	Server                       *ServerConfig          `json:"server,omitempty"`
	EndCallFunctionEnabled       bool                   `json:"endCallFunctionEnabled,omitempty"`
	EndCallMessage               *string                `json:"endCallMessage,omitempty"`
	EndCallPhrases               []string               `json:"endCallPhrases,omitempty"`
	VoicemailDetection           *VoicemailConfig       `json:"voicemailDetection,omitempty"`
	VoicemailMessage             *string                `json:"voicemailMessage,omitempty"`
	StartSpeakingPlan            *SpeakingPlan          `json:"startSpeakingPlan,omitempty"`
	AnalysisPlan                 *AnalysisPlan          `json:"analysisPlan,omitempty"`
	SilenceTimeoutSeconds        *int                   `json:"silenceTimeoutSeconds,omitempty"`
	MaxDurationSeconds           *int                   `json:"maxDurationSeconds,omitempty"`
	BackgroundSound              *string                `json:"backgroundSound,omitempty"`
	BackchannelingEnabled        *bool                  `json:"backchannelingEnabled,omitempty"`
	BackgroundDenoisingEnabled   *bool                  `json:"backgroundDenoisingEnabled,omitempty"`
	ModelOutputInMessagesEnabled *bool                  `json:"modelOutputInMessagesEnabled,omitempty"`
	VariableValues               map[string]any         `json:"variableValues,omitempty"`
}

// TranscriberConfig contains settings for speech-to-text
type TranscriberConfig struct {
	Model    string `json:"model"`
	Language string `json:"language"`
	Provider string `json:"provider"`
}

// VoicemailConfig contains settings for voicemail detection
type VoicemailConfig struct {
	Provider                           string   `json:"provider"`
	VoicemailDetectionTypes            []string `json:"voicemailDetectionTypes"`
	Enabled                            bool     `json:"enabled"`
	MachineDetectionTimeout            int      `json:"machineDetectionTimeout"`
	MachineDetectionSpeechThreshold    int      `json:"machineDetectionSpeechThreshold"`
	MachineDetectionSpeechEndThreshold int      `json:"machineDetectionSpeechEndThreshold"`
	MachineDetectionSilenceTimeout     int      `json:"machineDetectionSilenceTimeout"`
}

// TranscriptionEndpointingPlan contains timing settings for transcription
type TranscriptionEndpointingPlan struct {
	OnPunctuationSeconds   float64 `json:"onPunctuationSeconds"`
	OnNoPunctuationSeconds float64 `json:"onNoPunctuationSeconds"`
	OnNumberSeconds        float64 `json:"onNumberSeconds"`
}

// SpeakingPlan contains settings for speech timing
type SpeakingPlan struct {
	WaitSeconds float64 `json:"waitSeconds"`
}

// SchemaProperty represents a property in the structured data schema
type SchemaProperty struct {
	Type string `json:"type"`
}

// AnalysisPlan contains settings for call analysis
type AnalysisPlan struct {
	SummaryPlan         *SummaryPlan          `json:"summaryPlan,omitempty"`
	StructuredDataPlan  *StructuredDataPlan   `json:"structuredDataPlan,omitempty"`
	StructuredDataMulti []StructuredDataMulti `json:"structuredDataMulti,omitempty"`
}

type SummaryPlan struct {
	Messages       []Message `json:"messages"`
	TimeoutSeconds *int      `json:"timeoutSeconds,omitempty"` // defaults to 5 seconds
}

type StructuredDataPlan struct {
	Messages       []Message         `json:"messages"`
	Enabled        bool              `json:"enabled"` // defaults to false
	Schema         *minds.Definition `json:"schema,omitempty"`
	TimeoutSeconds *int              `json:"timeoutSeconds,omitempty"` // defaults to 5 seconds
}

type StructuredDataMulti struct {
	Key  string              `json:"key"`
	Plan *StructuredDataPlan `json:"plan,omitempty"`
}

type SuccessEvaluationPlan struct {
	// 	Rubric (enum) options include:
	// 		‘NumericScale’: A scale of 1 to 10.
	// 		‘DescriptiveScale’: A scale of Excellent, Good, Fair, Poor.
	// 		‘Checklist’: A checklist of criteria and their status.
	// 		‘Matrix’: A grid that evaluates multiple criteria across different performance levels.
	// 		‘PercentageScale’: A scale of 0% to 100%.
	// 		‘LikertScale’: A scale of Strongly Agree, Agree, Neutral, Disagree, Strongly Disagree.
	// 		‘AutomaticRubric’: Automatically break down evaluation into several criteria, each with its own score.
	// 		‘PassFail’: A simple 'true' if call passed, 'false' if not.
	//
	// Default is 'PassFail'.
	Rubric *string `json:"rubric,omitempty" enum:"NumericScale,DescriptiveScale,Checklist,Matrix,PercentageScale,LikertScale,AutomaticRubric,PassFail"`
	// 	These are the messages used to generate the success evaluation.
	//
	// @default: [ \{ "role": "system", "content": "You are an expert call
	// evaluator. You will be given a transcript of a call and the system prompt of
	// the AI participant. Determine if the call was successful based on the
	// objectives inferred from the system prompt. DO NOT return anything except the
	// result.\n\nRubric:\\n\{\{rubric}}\n\nOnly respond with the result." }, \{
	// "role": "user", "content": "Here is the transcript:\n\n\{\{transcript}}\n\n"
	// }, \{ "role": "user", "content": "Here was the system prompt of the
	// call:\n\n\{\{systemPrompt}}\n\n" } ]
	//
	// You can customize by providing any messages you want.
	//
	// Here are the template variables available:
	//
	// {{transcript}}: the transcript of the call from call.artifact.transcript-
	// {{systemPrompt}}: the system prompt of the call from
	// assistant.model.messages[type=system].content- {{rubric}}: the rubric of the
	// success evaluation from successEvaluationPlan.rubric
	Messages       []Message `json:"messages,omitempty"`
	Enabled        *bool     `json:"enabled,omitempty"`        // defaults to true
	TimeoutSeconds *int      `json:"timeoutSeconds,omitempty"` // defaults to 5 seconds
}

var DefaultElevenLabsVoiceConfig = ElevenLabsVoiceConfig{
	Model:                  "eleven_flash_v2_5",
	VoiceID:                "yM93hbw8Qtvdma2wCnJG",
	Provider:               "11labs",
	Stability:              0.5,
	SimilarityBoost:        0.75,
	FillerInjectionEnabled: false,
}

var DefaultModelConfig = ModelConfig{
	Provider:    "openai",
	Model:       "gpt-4o-mini",
	MaxTokens:   300,
	Temperature: 1.0,
}

var DefaultTranscriber = TranscriberConfig{
	Model:    "nova-2",
	Language: "en",
	Provider: "deepgram",
}

func DefaultAssistant(agentName, prompt, firstMessage, webhook, voicemailMessage, endCallMessage string) (*Assistant, error) {
	// VoicemailConfig{
	// 	Provider: "twilio",
	// 	VoicemailDetectionTypes: []string{
	// 		"machine_start",
	// 		"machine_end_beep",
	// 		"machine_end_silence",
	// 		"machine_end_other",
	// 	},
	// 	Enabled:                            true,
	// 	MachineDetectionTimeout:            25,
	// 	MachineDetectionSpeechThreshold:    2500,
	// 	MachineDetectionSpeechEndThreshold: 2000,
	// 	MachineDetectionSilenceTimeout:     2000,
	// },

	req := &Assistant{
		Name:                   &agentName,
		Voice:                  &DefaultElevenLabsVoiceConfig,
		Model:                  &DefaultModelConfig,
		Transcriber:            &DefaultTranscriber,
		FirstMessage:           &firstMessage,
		ClientMessages:         []string{},
		ServerMessages:         []string{"end-of-call-report", "function-call", "tool-calls"},
		Server:                 nil,
		EndCallFunctionEnabled: true,
		EndCallMessage:         &endCallMessage,
		EndCallPhrases:         []string{"goodbye"},
		VoicemailDetection:     nil,
		VoicemailMessage:       &voicemailMessage,
		AnalysisPlan:           nil,
		// SilenceTimeoutSeconds:        30,
		// MaxDurationSeconds:           180,
		// BackgroundSound:              "off",
		// BackchannelingEnabled:        false,
		// BackgroundDenoisingEnabled:   false,
		// ModelOutputInMessagesEnabled: false,
	}

	return req, nil
}

// GetAssistant retrieves an assistant by its ID and saves the raw JSON response to a file
func GetAssistant(ctx context.Context, id string) (*Assistant, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.vapi.ai/assistant/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for assistant: %w", err)
	}

	apiKey := os.Getenv("VAPI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("VAPI_API_KEY not set")
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to http client failed: %w", err)
	}
	defer resp.Body.Close()

	var body bytes.Buffer
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get assistant. code: %d msg: %s", resp.StatusCode, body.String())
	}

	// Save the raw JSON response to a file for examination
	filename := fmt.Sprintf("assistant-%s-response.json", id)
	if err := os.WriteFile(filename, body.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("failed to save response to file: %w", err)
	}

	var result Assistant
	if err := json.Unmarshal(body.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
