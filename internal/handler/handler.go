package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"fsocket/internal/hub"
	"fsocket/pkg/response"
)

var messagesPublished int64

type PublishRequest struct {
	Channel  string          `json:"channel"`
	EventType string         `json:"eventType"`
	Message  string          `json:"msg"`
	Extra    json.RawMessage `json:"extra,omitempty"`
}

func SSE(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := r.URL.Query().Get("channel")
		if channel == "" {
			http.Error(w, "channel is required", http.StatusBadRequest)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		client := h.Subscribe(channel)
		defer h.Unsubscribe(client)

	notifyLoop:
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					break notifyLoop
				}

				var payload map[string]interface{}
				if err := json.Unmarshal(msg, &payload); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}

				eventType := "message"
				if et, ok := payload["eventType"].(string); ok {
					eventType = et
				}

				fmt.Fprintf(w, "event: %s\n", eventType)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()

			case <-r.Context().Done():
				break notifyLoop
			}
		}
	}
}

func Publish(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		defer r.Body.Close()

		var req PublishRequest
		if err := json.Unmarshal(body, &req); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		if req.Channel == "" || req.Message == "" {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "channel and message are required"})
			return
		}

		eventType := req.EventType
		if eventType == "" {
			eventType = "message"
		}

		msg := &hub.Message{
			EventType: eventType,
			Message:   req.Message,
			Extra:     req.Extra,
			Timestamp: "timestamp",
		}

		extraMap := make(map[string]interface{})
		if req.Extra != nil {
			json.Unmarshal(req.Extra, &extraMap)
		}
		extraMap["channel"] = req.Channel
		if msg.Timestamp == "timestamp" {
			extraMap["timestamp"] = "timestamp"
		}

		extraBytes, _ := json.Marshal(extraMap)
		msg.Extra = extraBytes

		h.Publish(req.Channel, msg)
		atomic.AddInt64(&messagesPublished, 1)

		response.JSON(w, http.StatusOK, map[string]string{"status": "message published"})
	}
}

func Broadcast(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		defer r.Body.Close()

		var req PublishRequest
		if err := json.Unmarshal(body, &req); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		eventType := req.EventType
		if eventType == "" {
			eventType = "message"
		}

		msg := &hub.Message{
			EventType: eventType,
			Message:   req.Message,
			Extra:     req.Extra,
		}

		h.Broadcast(msg)
		atomic.AddInt64(&messagesPublished, 1)

		response.JSON(w, http.StatusOK, map[string]string{"status": "message broadcasted"})
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Stats(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channels, clients := h.GetStats()
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"activeChannels":     channels,
			"connectedClients":   clients,
			"messagesPublished": atomic.LoadInt64(&messagesPublished),
		})
	}
}
