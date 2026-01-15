package httppack

import (
	"copr/worker"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type HTTPHandlers struct {
	workers *worker.Porkers
	conn    *pgx.Conn
}

func NewHTTPHandlers(works *worker.Porkers, conn *pgx.Conn) *HTTPHandlers {
	return &HTTPHandlers{
		workers: works,
		conn:    conn,
	}
}

func (h *HTTPHandlers) CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var workDTO worker.Worker
	if err := json.NewDecoder(r.Body).Decode(&workDTO); err != nil {
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errDto)
		return
	}
	defer r.Body.Close()

	if workDTO.FullName == "" || workDTO.Position == "" {
		errDto := ErrorDTO{
			Message: "FullName and Position are required fields",
			Time:    time.Now(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errDto)
		return
	}

	query := `
        INSERT INTO personal (full_name, position)
        VALUES ($1, $2)
        RETURNING id;
    `

	var insertedId int
	err := h.conn.QueryRow(r.Context(), query,
		workDTO.FullName,
		workDTO.Position,
	).Scan(&insertedId)

	if err != nil {
		log.Printf("Database insert error: %v", err)

		errDto := ErrorDTO{
			Message: "Failed to insert into database: " + err.Error(),
			Time:    time.Now(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errDto)
		return
	}

	personal := worker.NewWorker(workDTO.Id, workDTO.FullName, workDTO.Position)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        insertedId,
		"full_name": workDTO.FullName,
		"position":  workDTO.Position,
		"message":   "Employee created successfully",
		"data":      personal,
	})
}

func (h *HTTPHandlers) GetWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	querySql := `
	SELECT id, full_name, position 
	FROM personal
	ORDER BY id ASC
	`
	rows, err := h.conn.Query(r.Context(), querySql)
	if err != nil {
		log.Printf("Database query error: %v", err)
		errDto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errDto)
		return
	}
	defer rows.Close()

	personal := make([]worker.Worker, 0)

	for rows.Next() {
		var prsl worker.Worker

		err := rows.Scan(
			&prsl.Id,
			&prsl.FullName,
			&prsl.Position,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			errDto := ErrorDTO{
				Message: "Failed to scan row: " + err.Error(),
				Time:    time.Now(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errDto)
			return
		}
		personal = append(personal, prsl)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		errDto := ErrorDTO{
			Message: "Database error: " + err.Error(),
			Time:    time.Now(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errDto)
		return
	}
	if len(personal) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(personal)

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

	sqlQuery := `
	DELETE FROM personal
	WHERE id = $1;
	`
	result, err := h.conn.Exec(r.Context(), sqlQuery, id)
	if err != nil {
		log.Printf("Database delete error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete from database: " + err.Error(),
		})
		return
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Worker not found with ID: " + idStr,
		})
		return
	}

	log.Printf("Worker with ID %d deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}
