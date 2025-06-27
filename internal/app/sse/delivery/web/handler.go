package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/messagebroker"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/utils"
)

type SseHandler struct {
	SseHandlerOptions
}

type SseHandlerOptions struct {
	logger                *slog.Logger
	inMemoryMessageBroker *messagebroker.InMemoryMessageBroker
}

func NewSseHandler(options SseHandlerOptions) SseHandler {
	return SseHandler{options}
}

func (s SseHandler) GetRewards(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	err := rc.SetWriteDeadline(time.Time{})
	if err != nil {
		s.logger.Error("unsupported .SetWriteDeadline()", slog.String("error", err.Error()))
		utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	ctx := r.Context()
	userId, ok := requestctx.UserId(ctx)
	if !ok {
		s.logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.logger.Error("unsupported streaming using http.Flusher")
		utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
	}

	subscriber := s.inMemoryMessageBroker.Subscribe(strconv.Itoa(userId))
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("client disconnected", slog.Int("userId", userId))
			return
		case msg := <-subscriber:
			// Filtering the stream by 'rewards' event name
			fmt.Fprintf(w, "event: rewards\n")
			// This is a specific format for Server-Sent Events (SSE) response
			fmt.Fprintf(w, "data: {\"points\": %d}\n\n", msg.Points)
			flusher.Flush()
		}
	}
}
