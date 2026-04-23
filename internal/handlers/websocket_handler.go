package handlers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"
	internalws "anonymous-communication/backend/internal/websocket"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WebSocketHandler struct {
	chatService *services.ChatService
	hub         *internalws.Hub
	jwtConfig   config.JWTConfig
	corsConfig  config.CORSConfig
	rateLimiter *middleware.RateLimiter
}

func NewWebSocketHandler(chatService *services.ChatService, hub *internalws.Hub, jwtConfig config.JWTConfig, corsConfig config.CORSConfig, rateLimiter *middleware.RateLimiter) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
		hub:         hub,
		jwtConfig:   jwtConfig,
		corsConfig:  corsConfig,
		rateLimiter: rateLimiter,
	}
}

func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	if !fiberws.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Verify origin before upgrading
	origin := c.Get(fiber.HeaderOrigin)
	if !h.IsOriginAllowed(origin) {
		return utils.Error(c, fiber.StatusForbidden, "origin not allowed", nil)
	}

	// Try to get token from cookie first
	token := strings.TrimSpace(c.Cookies(h.jwtConfig.CookieName))

	// If no token in cookie, try query parameter
	if token == "" {
		token = strings.TrimSpace(c.Query("token"))
	}

	// If still no token, try Authorization header
	if token == "" {
		authHeader := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:] // Extract token after "Bearer "
		}
	}

	if token == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	claims, err := utils.ParseToken(h.jwtConfig.Secret, token)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "invalid or expired session", nil)
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}

func (h *WebSocketHandler) HandleConnection(conn *fiberws.Conn) {
	userIDValue := conn.Locals("user_id")
	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		_ = conn.WriteJSON(internalws.NewErrorEvent("authentication required"))
		_ = conn.Close()
		return
	}

	client := internalws.NewClient(userID, conn)
	h.hub.Register(client)
	defer func() {
		h.hub.Unregister(client)
		_ = client.Close()
	}()

	_ = client.WriteJSON(internalws.NewConnectedEvent())

	senderID, err := utils.ParseUUID(userID)
	if err != nil {
		_ = client.WriteJSON(internalws.NewErrorEvent("invalid authenticated user"))
		return
	}

	for {
		var inbound internalws.IncomingMessage
		if err := conn.ReadJSON(&inbound); err != nil {
			break
		}

		messageType := inbound.Type
		if messageType == "" {
			messageType = internalws.EventTypeMessage
		}

		if messageType == internalws.EventTypeMessage {
			h.handleChatWebSocketMessage(client, userID, senderID, inbound)
			continue
		}

		if isCallWebSocketEventType(messageType) {
			h.handleCallWebSocketMessage(client, userID, inbound)
			continue
		}

		_ = client.WriteJSON(internalws.NewErrorEvent("unsupported websocket event type"))
	}
}

func (h *WebSocketHandler) handleChatWebSocketMessage(client *internalws.Client, userID string, senderID uuid.UUID, inbound internalws.IncomingMessage) {
	if h.rateLimiter != nil {
		allowed, retryAfter := h.rateLimiter.AllowMessageForUser(userID)
		if !allowed {
			retryAfterSeconds := int(math.Ceil(retryAfter.Seconds()))
			if retryAfterSeconds < 1 {
				retryAfterSeconds = 1
			}

			_ = client.WriteJSON(internalws.NewErrorEvent(fmt.Sprintf("message rate limit exceeded. try again in %d seconds.", retryAfterSeconds)))
			return
		}
	}

	message, err := h.chatService.SendMessage(context.Background(), senderID, models.SendMessageRequest{
		ReceiverID: inbound.ReceiverID,
		Content:    inbound.Content,
	})
	if err != nil {
		_ = client.WriteJSON(internalws.NewErrorEvent(h.websocketErrorMessage(err)))
		return
	}

	event := internalws.NewMessageEvent(message)
	_ = client.WriteJSON(event)
	internalws.BroadcastToUser(h.hub, message.ReceiverID, event)
}

func (h *WebSocketHandler) handleCallWebSocketMessage(client *internalws.Client, userID string, inbound internalws.IncomingMessage) {
	receiverID := strings.TrimSpace(inbound.ReceiverID)
	callID := strings.TrimSpace(inbound.CallID)

	if receiverID == "" {
		_ = client.WriteJSON(internalws.NewErrorEvent("receiver_id is required for calls"))
		return
	}

	if _, err := utils.ParseUUID(receiverID); err != nil {
		_ = client.WriteJSON(internalws.NewErrorEvent("receiver_id must be a valid UUID"))
		return
	}

	if receiverID == userID {
		_ = client.WriteJSON(internalws.NewErrorEvent("cannot call yourself"))
		return
	}

	if callID == "" {
		_ = client.WriteJSON(internalws.NewErrorEvent("call_id is required for call events"))
		return
	}

	switch inbound.Type {
	case internalws.EventTypeCallOffer, internalws.EventTypeCallAnswer:
		if inbound.Description == nil || strings.TrimSpace(inbound.Description.Type) == "" || strings.TrimSpace(inbound.Description.SDP) == "" {
			_ = client.WriteJSON(internalws.NewErrorEvent("description is required for this call event"))
			return
		}
	case internalws.EventTypeCallICECandidate:
		if inbound.Candidate == nil || strings.TrimSpace(inbound.Candidate.Candidate) == "" {
			_ = client.WriteJSON(internalws.NewErrorEvent("candidate is required for ICE events"))
			return
		}
	case internalws.EventTypeCallMuteState:
		if inbound.Muted == nil {
			_ = client.WriteJSON(internalws.NewErrorEvent("muted state is required"))
			return
		}
	case internalws.EventTypeCallVideoState:
		if inbound.VideoOff == nil {
			_ = client.WriteJSON(internalws.NewErrorEvent("video state is required"))
			return
		}
	}

	event := internalws.NewCallEvent(
		inbound.Type,
		callID,
		userID,
		strings.TrimSpace(inbound.Username),
		strings.TrimSpace(inbound.Reason),
		inbound.Description,
		inbound.Candidate,
		inbound.Muted,
		inbound.VideoOff,
		inbound.CallType,
	)

	delivered := internalws.BroadcastToUser(h.hub, receiverID, event)
	if delivered {
		return
	}

	switch inbound.Type {
	case internalws.EventTypeCallCancel, internalws.EventTypeCallDecline, internalws.EventTypeCallEnd:
		return
	default:
		_ = client.WriteJSON(internalws.NewErrorEvent("call participant is offline or unavailable"))
	}
}

func isCallWebSocketEventType(messageType string) bool {
	switch messageType {
	case internalws.EventTypeCallInvite,
		internalws.EventTypeCallAccept,
		internalws.EventTypeCallDecline,
		internalws.EventTypeCallBusy,
		internalws.EventTypeCallCancel,
		internalws.EventTypeCallOffer,
		internalws.EventTypeCallAnswer,
		internalws.EventTypeCallICECandidate,
		internalws.EventTypeCallEnd,
		internalws.EventTypeCallMuteState,
		internalws.EventTypeCallVideoState:
		return true
	default:
		return false
	}
}

func (h *WebSocketHandler) IsOriginAllowed(origin string) bool {
	// Parse comma-separated origins from config.
	// In a real app, you might want to cache this split slice.
	allowedOrigins := strings.Split(h.corsConfig.AllowedOrigin, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"ghostline.reporoot.in", "localhost:3000"} 
	}

	return middleware.IsOriginAllowed(origin, allowedOrigins)
}

func (h *WebSocketHandler) websocketErrorMessage(err error) string {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		for _, message := range validationErr.Fields {
			return message
		}
		return validationErr.Error()
	case errors.Is(err, models.ErrUserNotFound):
		return "receiver not found"
	case errors.Is(err, models.ErrCannotMessageSelf):
		return err.Error()
	default:
		return "unable to process websocket message"
	}
}
