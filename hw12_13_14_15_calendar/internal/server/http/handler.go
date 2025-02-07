package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	app    Application
	logger Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	var req CreateRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	ctx := context.Background()

	err = h.app.CreateEvent(ctx, req.ID, req.Title)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	renderSuccessResponse(w, Response{})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	var req UpdateEventRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	ctx := context.Background()

	event := app.Event{
		ID:          req.ID,
		Title:       req.Title,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
		UserID:      req.UserID,
		NotifyDays:  req.NotifyDays,
	}

	err = h.app.UpdateEvent(ctx, id, event)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	renderSuccessResponse(w, Response{})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]

	ctx := context.Background()

	err := h.app.DeleteEvent(ctx, eventID)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	renderSuccessResponse(w, Response{})
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	searchDate := params.Get("date")
	searchPeriod := params.Get("period")

	ctx := context.Background()

	events, err := h.app.ListEvents(ctx, searchDate, searchPeriod)
	if err != nil {
		h.logger.Error(err.Error())
		renderErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	resp := EventsResponse{
		Events: make([]Event, 0, len(events)),
	}

	for _, event := range events {
		resp.Events = append(resp.Events, Event{
			ID:          event.ID,
			Title:       event.Title,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Description: event.Description,
			UserID:      event.UserID,
			NotifyDays:  event.NotifyDays,
		})
	}

	renderSuccessResponse(w, resp)
}

func renderSuccessResponse(w http.ResponseWriter, resp interface{}) {
	jsonResp, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func renderErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	resp := Response{
		Error:   true,
		Message: err.Error(),
	}

	jsonResp, _ := json.Marshal(resp)

	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}
