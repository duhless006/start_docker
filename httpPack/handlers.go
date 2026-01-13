package httppack

import (
	"copr/worker"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	workers *worker.Porkers
}

func NewHTTPHandlers(works *worker.Porkers) *HTTPHandlers {
	return &HTTPHandlers{
		workers: works,
	}
}

func (h *HTTPHandlers) CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	var workDTO worker.Worker
	if err := json.NewDecoder(r.Body).Decode(&workDTO); err != nil {
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDto.ToString(), http.StatusBadRequest)
		return
	}
	personal := worker.NewWorker(workDTO.Id, workDTO.FullName, workDTO.Position)
	if err := h.workers.AddWorker(*personal); err != nil {
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, worker.ErrWorkerAlreadyExists) {
			http.Error(w, errDto.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDto.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(personal, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

func (h *HTTPHandlers) GetWorkerHandler(w http.ResponseWriter, r *http.Request) {
	person, err := h.workers.GetWorker()
	if err != nil {
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, worker.ErrWorkerNotFound) {
			http.Error(w, errDto.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDto.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

func (h *HTTPHandlers) DeleteWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(r)["id"]
	log.Printf("Extracted ID from URL: %s", idStr)
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Параметр id обязателен"})
		return

	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting ID '%s' to int: %v", idStr, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID должен быть числом"})
		return
	}

	if err := h.workers.DeleteWorker(id); err != nil {
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, worker.ErrWorkerNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(errDto)
		return

	}
	log.Printf("Worker with ID %d deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}
