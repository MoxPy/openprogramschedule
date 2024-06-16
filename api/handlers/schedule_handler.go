package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"openprogramschedule/internal/models"
	"openprogramschedule/internal/repository"
	"openprogramschedule/internal/validators"
	"strconv"
)

type ScheduleHandler struct {
	Db *sql.DB
}

func (env *ScheduleHandler) AddScheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var scheduleData models.Schedule

		err := json.NewDecoder(r.Body).Decode(&scheduleData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		if err = validators.ValidateSchedule(&scheduleData); err != nil {
			http.Error(w, fmt.Sprintf("Validation Error: %v", err), http.StatusBadRequest)
			return
		}

		id, err := repository.AddSchedule(&scheduleData, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf("Added new schedule with id: %v", id)
		response := map[string]interface{}{
			"id":      id,
			"message": message,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) GetAllSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		schedules, err := repository.GetAllSchedules(env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}
		if len(schedules) == 0 {
			log.Printf("No results found")
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(schedules)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) GetScheduleByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing schedule ID", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
			return
		}

		program, err := repository.GetScheduleByID(id, env.Db)
		if err != nil {
			if err.Error() == "schedule not found" {
				http.Error(w, "Schedule not found: invalid ID", http.StatusNotFound)
				return
			}
			log.Printf("Error during operation: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(program)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) GetScheduleByProgramIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		programIdStr := r.URL.Query().Get("programId")
		if programIdStr == "" {
			http.Error(w, "Missing program id", http.StatusBadRequest)
			return
		}
		programIdInt, err := strconv.Atoi(programIdStr)
		programId := uint(programIdInt)
		if err != nil {
			http.Error(w, "Invalid program ID", http.StatusBadRequest)
			return
		}
		schedules, err := repository.GetScheduleByProgramID(programId, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}
		if len(*schedules) == 0 {
			log.Printf("No results found")
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(schedules)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) GetScheduleByDayHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dayStr := r.URL.Query().Get("day")
		if dayStr == "" {
			http.Error(w, "Missing day", http.StatusBadRequest)
			return
		}
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			http.Error(w, "Invalid day", http.StatusBadRequest)
			return
		}
		if day > 7 || day < 1 {
			http.Error(w, "Day must be between 1 and 7", http.StatusBadRequest)
			return
		}
		schedules, err := repository.GetScheduleByDay(day, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}
		if len(*schedules) == 0 {
			log.Printf("No schedules found for day: %v", day)
			http.Error(w, "No schedule found for this day", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(schedules)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) GetScheduleByDateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dayStr := r.URL.Query().Get("date")
		if dayStr == "" {
			http.Error(w, "Missing day", http.StatusBadRequest)
			return
		}

		schedules, err := repository.GetScheduleByDate(dayStr, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			return
		}
		if len(*schedules) == 0 {
			log.Printf("No schedules found for day: %v", dayStr)
			http.Error(w, "No schedule found for this day", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(schedules)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) UpdateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
			return
		}

		var updatedSchedule models.Schedule
		err = json.NewDecoder(r.Body).Decode(&updatedSchedule)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing body: %v", err)
			}
		}(r.Body)

		if err = validators.ValidateSchedule(&updatedSchedule); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repository.UpdateScheduleByID(id, updatedSchedule, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		updatedSchedule.Id = new(uint)
		*updatedSchedule.Id = id
		response := map[string]interface{}{
			"program": updatedSchedule,
			"message": "Update successful",
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ScheduleHandler) DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid program ID", http.StatusBadRequest)
			return
		}
		err = repository.DeleteScheduleByID(id, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			return
		}
		message := fmt.Sprintf("Deleted schedule: %v", id)
		response := map[string]interface{}{
			"message": message,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func (env *ScheduleHandler) DeleteAllSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		err := repository.DeleteAllSchedules(env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"message": "Deleted all schedules",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
