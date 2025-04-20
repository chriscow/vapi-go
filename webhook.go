package vapi

import "time"

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

/*
type
"assistant-request"
Required
This is the type of the message. “assistant-request” is sent to fetch assistant configuration for an incoming call.

phoneNumber
object
Optional
This is the phone number associated with the call.

This matches one of the following:

call.phoneNumber,
call.phoneNumberId.

Show 4 variants
timestamp
double
Optional
This is the timestamp of when the message was sent in milliseconds since Unix Epoch.

artifact
object
Optional
This is a live version of the call.artifact.

This matches what is stored on call.artifact after the call.

Show 8 properties
assistant
object
Optional
This is the assistant that is currently active. This is provided for convenience.

This matches one of the following:

call.assistant,
call.assistantId,
call.squad[n].assistant,
call.squad[n].assistantId,
call.squadId->[n].assistant,
call.squadId->[n].assistantId.

Show 33 properties
customer
object
Optional
This is the customer associated with the call.

This matches one of the following:

call.customer,
call.customerId.

Show 6 properties
call
object
Optional
This is the call object.

This matches what was returned in POST /call.

Note: This might get stale during the call. To get the latest call object, especially after the call is ended, use GET /call/:id.
*/

/*
"map[message:map[call:map[assistantId:<nil> assistantOverrides:map[variableValues:map[account-sid:c033b672-5b99-42ae-9ce2-b231e9a522fb application-sid:79d078c8-76b2-452a-99e0-ddd5abbf6269 cid:2270729304-3954070205-1530341129@msc1.382COM.COM forwarded-for:64.125.111.10 originating-carrier:382com voip-carrier-sid:a5569621-84ac-49cc-a8b8-11c7fb96b905]] createdAt:2025-04-19T16:50:05.656Z customer:map[number:+14258295443 sipUri:sip:+14258295443@44.229.228.186:5060] id:0ffa0790-6ad0-4b43-ba0f-7b665080f4dc orgId:4ab76b87-4964-48a6-8d82-1b32ec845652 phoneCallProvider:vapi phoneCallProviderDetails:map[sbcCallId:2270729304-3954070205-1530341129@msc1.382COM.COM] phoneCallProviderId:2e9501c3-9da0-446e-bf5a-bafe27e045fb phoneCallTransport:sip phoneNumberId:fdaffae3-edf4-4969-a273-690e04af6135 squadId:<nil> status:ringing type:inboundPhoneCall updatedAt:2025-04-19T16:50:05.656Z] customer:map[number:+14258295443 sipUri:sip:+14258295443@44.229.228.186:5060] phoneNumber:map[assistantId:<nil> authentication:<nil> createdAt:2025-04-19T16:39:11.998Z credentialId:<nil> credentialIds:<nil> fallbackDestination:<nil> fallbackForwardingPhoneNumber:<nil> hooks:<nil> id:fdaffae3-edf4-4969-a273-690e04af6135 name:Unlabeled number:+16592555076 numberE164CheckEnabled:<nil> orgId:4ab76b87-4964-48a6-8d82-1b32ec845652 provider:vapi providerResourceId:45719bef-ac51-496b-b8e2-1241bf29e128 server:<nil> serverUrl:<nil> serverUrlSecret:<nil> sipUri:<nil> squadId:<nil> status:active stripeSubscriptionCurrentPeriodStart:<nil> stripeSubscriptionId:<nil> stripeSubscriptionStatus:<nil> twilioAccountSid:<nil> twilioAuthToken:<nil> twilioOutgoingCallerId:<nil> updatedAt:2025-04-19T16:41:12.240Z useClusterSip:<nil>] timestamp:1.745081405664e+12 type:assistant-request]]"
*/

type AssistantRequestEnvelope struct {
	AssistantRequest struct {
		Type        string       `json:"type"`
		PhoneNumber *PhoneNumber `json:"phoneNumber,omitempty"`
		Timestamp   *float64     `json:"timestamp,omitempty"`
		Artifact    *Artifact    `json:"artifact,omitempty"`
		Assistant   *Assistant   `json:"assistant,omitempty"`
		Customer    *Customer    `json:"customer,omitempty"`
		Call        *Call        `json:"call,omitempty"`
	} `json:"message"`
}

/*
Server Message Response Assistant Request
object

Hide 7 properties
destination
object
Optional
This is the destination to transfer the inbound call to. This will immediately transfer without using any assistants.

If this is sent, assistantId, assistant, squadId, and squad are ignored.

Show 2 variants
assistantId
string
Optional
This is the assistant that will be used for the call. To use a transient assistant, use assistant instead.

assistant
object
Optional
This is the assistant that will be used for the call. To use an existing assistant, use assistantId instead.

If you’re unsure why you’re getting an invalid assistant, try logging your response and send the JSON blob to POST /assistant which will return the validation errors.

Show 33 properties
assistantOverrides
object
Optional
These are the overrides for the assistant or assistantId’s settings and template variables.

Show 34 properties
squadId
string
Optional
This is the squad that will be used for the call. To use a transient squad, use squad instead.

squad
object
Optional
This is a squad that will be used for the call. To use an existing squad, use squadId instead.

Show 3 properties
error
string
Optional
This is the error if the call shouldn’t be accepted. This is spoken to the customer.

If this is sent, assistantId, assistant, squadId, squad, and destination are ignored.
*/
type AssistantRequestResponse struct {
	Destination        *Destination `json:"destination,omitempty"`
	AssistantId        *string      `json:"assistantId,omitempty"`
	Assistant          *Assistant   `json:"assistant,omitempty"`
	AssistantOverrides *Assistant   `json:"assistantOverrides,omitempty"`
}

// EndOfCallReportEnvelope represents the report generated at the end of a call
type EndOfCallReportEnvelope struct {
	Report struct {
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
	} `json:"message"`
}
