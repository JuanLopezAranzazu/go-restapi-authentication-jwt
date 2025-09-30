package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/db"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/middlewares"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/models"
	"github.com/gorilla/mux"
)

type CreateEventRequest struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
}

type EventResponse struct {
	ID     uint      `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserID uint      `json:"user_id"`
}

// crear eventos
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	event := models.Event{
		Title:  req.Title,
		Date:   req.Date,
		UserID: userID,
	}

	if err := db.DB.Create(&event).Error; err != nil {
		http.Error(w, "No se pudo crear el evento", http.StatusInternalServerError)
		return
	}

	resp := EventResponse{
		ID:     event.ID,
		Title:  event.Title,
		Date:   event.Date,
		UserID: event.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// obtener eventos del usuario
func GetMyEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// buscar eventos
	var events []models.Event
	if err := db.DB.Where("user_id = ?", userID).Find(&events).Error; err != nil {
		http.Error(w, "No se pudieron obtener los eventos", http.StatusInternalServerError)
		return
	}

	resp := []EventResponse{}
	for _, e := range events {
		resp = append(resp, EventResponse{
			ID:     e.ID,
			Title:  e.Title,
			Date:   e.Date,
			UserID: e.UserID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// obtener un evento especifico
func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParam)
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// buscar evento
	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		http.Error(w, "Evento no encontrado", http.StatusNotFound)
		return
	}

	// validar el usuario
	if event.UserID != userID {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	resp := EventResponse{
		ID:     event.ID,
		Title:  event.Title,
		Date:   event.Date,
		UserID: event.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// actualizar evento
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParam)
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// buscar evento
	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		http.Error(w, "Evento no encontrado", http.StatusNotFound)
		return
	}

	// validar el usuario
	if event.UserID != userID {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event.Title = req.Title
	event.Date = req.Date
	event.UpdatedAt = time.Now()

	if err := db.DB.Save(&event).Error; err != nil {
		http.Error(w, "No se pudo actualizar el evento", http.StatusInternalServerError)
		return
	}

	resp := EventResponse{
		ID:     event.ID,
		Title:  event.Title,
		Date:   event.Date,
		UserID: event.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// eliminar un evento
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParam)
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// buscar evento
	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		http.Error(w, "Evento no encontrado", http.StatusNotFound)
		return
	}

	// validar el usuario
	if event.UserID != userID {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	if err := db.DB.Delete(&event).Error; err != nil {
		http.Error(w, "No se pudo eliminar el evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
