package websocket

import "anonymous-communication/backend/internal/models"

const (
	EventTypeConnected        = "connected"
	EventTypeMessage          = "message"
	EventTypeError            = "error"
	EventTypeCallInvite       = "call_invite"
	EventTypeCallAccept       = "call_accept"
	EventTypeCallDecline      = "call_decline"
	EventTypeCallBusy         = "call_busy"
	EventTypeCallCancel       = "call_cancel"
	EventTypeCallOffer        = "call_offer"
	EventTypeCallAnswer       = "call_answer"
	EventTypeCallICECandidate = "call_ice_candidate"
	EventTypeCallEnd          = "call_end"
	EventTypeCallMuteState    = "call_mute_state"
)

type SessionDescriptionPayload struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type ICECandidatePayload struct {
	Candidate        string  `json:"candidate"`
	SDPMid           *string `json:"sdpMid,omitempty"`
	SDPMLineIndex    *uint16 `json:"sdpMLineIndex,omitempty"`
	UsernameFragment *string `json:"usernameFragment,omitempty"`
}

type IncomingMessage struct {
	Type        string                     `json:"type"`
	ReceiverID  string                     `json:"receiver_id"`
	Content     string                     `json:"content"`
	CallID      string                     `json:"call_id,omitempty"`
	Description *SessionDescriptionPayload `json:"description,omitempty"`
	Candidate   *ICECandidatePayload       `json:"candidate,omitempty"`
	Reason      string                     `json:"reason,omitempty"`
	Username    string                     `json:"username,omitempty"`
	Muted       *bool                      `json:"muted,omitempty"`
}

type Event struct {
	Type        string                     `json:"type"`
	Message     *models.MessageResponse    `json:"message,omitempty"`
	Error       string                     `json:"error,omitempty"`
	CallID      string                     `json:"call_id,omitempty"`
	UserID      string                     `json:"user_id,omitempty"`
	Username    string                     `json:"username,omitempty"`
	Description *SessionDescriptionPayload `json:"description,omitempty"`
	Candidate   *ICECandidatePayload       `json:"candidate,omitempty"`
	Reason      string                     `json:"reason,omitempty"`
	Muted       *bool                      `json:"muted,omitempty"`
}

func NewConnectedEvent() Event {
	return Event{Type: EventTypeConnected}
}

func NewMessageEvent(message *models.MessageResponse) Event {
	return Event{
		Type:    EventTypeMessage,
		Message: message,
	}
}

func NewErrorEvent(message string) Event {
	return Event{
		Type:  EventTypeError,
		Error: message,
	}
}

func NewCallEvent(eventType, callID, userID, username, reason string, description *SessionDescriptionPayload, candidate *ICECandidatePayload, muted *bool) Event {
	return Event{
		Type:        eventType,
		CallID:      callID,
		UserID:      userID,
		Username:    username,
		Reason:      reason,
		Description: description,
		Candidate:   candidate,
		Muted:       muted,
	}
}
