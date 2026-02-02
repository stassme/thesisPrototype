// handler does HTTP only: parse request, call service, write response
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"thesisPrototype/internal/logging"
	"thesisPrototype/internal/service"
)

// Handler holds the processor and logger; main wires these in
type Handler struct {
	Processor service.Processor
	Logger    *slog.Logger
	Timeout   time.Duration
}

// Health returns 200 and {"status":"ok"} for liveness checks
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Process accepts GET (default payload) or POST with JSON body, calls service -> returns JSON
func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "none"
	}
	log := h.Logger
	logging.RequestStart(log, r.Method, r.URL.Path, requestID)
	start := time.Now()
	status := http.StatusOK
	defer func() {
		logging.RequestEnd(log, r.Method, r.URL.Path, requestID, status, time.Since(start).Milliseconds())
	}()

	ctx, cancel := context.WithTimeout(r.Context(), h.Timeout)
	defer cancel()

	var input service.ProcessInput
	if r.Method == http.MethodPost && r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logging.LogError(log, err, "decode request body", "request_id", requestID)
			status = http.StatusBadRequest
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"invalid json"}`))
			return
		}
	}
	if input.Payload == "" {
		input.Payload = "hello"
	}

	result, err := h.Processor.Process(ctx, input)
	if err != nil {
		if ctx.Err() != nil {
			logging.LogError(log, err, "request cancelled or timeout", "request_id", requestID)
			status = http.StatusGatewayTimeout
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"timeout"}`))
		} else {
			logging.LogError(log, err, "process failed", "request_id", requestID)
			status = http.StatusInternalServerError
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"internal"}`))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(result)
}
