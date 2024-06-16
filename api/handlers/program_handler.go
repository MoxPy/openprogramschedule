package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"openprogramschedule/internal/models"
	"openprogramschedule/internal/repository"
	"openprogramschedule/internal/validators"
	"strconv"
)

type ProgramHandler struct {
	Db *sql.DB
}

func (env *ProgramHandler) AddProgramHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var programData models.Program

		err := json.NewDecoder(r.Body).Decode(&programData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		if err = validators.ValidateProgram(&programData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := repository.AddProgram(&programData, env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf("Added new program with id: %v", id)
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

func (env *ProgramHandler) GetProgramByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing program ID", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid program ID", http.StatusBadRequest)
			return
		}

		program, err := repository.GetProgramByID(id, env.Db)
		if err != nil {
			if err.Error() == "program not found" {
				http.Error(w, "Program not found: invalid ID", http.StatusNotFound)
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

func (env *ProgramHandler) GetProgramByNameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing program name", http.StatusBadRequest)
			return
		}

		program, err := repository.GetProgramByName(name, env.Db)
		if err != nil {
			if err.Error() == "program not found" {
				http.Error(w, "Program not found", http.StatusNotFound)
				return
			}
			log.Printf("Error during program retrieval: %v", err)
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

	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func (env *ProgramHandler) GetProgramsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		category := r.URL.Query().Get("category")
		if category == "" {
			http.Error(w, "Missing category parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Received category: '%s'", category)
		programs, err := repository.GetProgramsByCategory(category, env.Db)
		if err != nil {
			log.Printf("Error during programs retrieval: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(programs)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func (env *ProgramHandler) GetAllProgramsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		programs, err := repository.GetAllPrograms(env.Db)
		if err != nil {
			log.Printf("Error during operation: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(programs) == 0 {
			http.Error(w, "No programs found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(programs)
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

func (env *ProgramHandler) UpdateProgramHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing program ID", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid program ID", http.StatusBadRequest)
			return
		}

		var updatedProgram models.Program
		err = json.NewDecoder(r.Body).Decode(&updatedProgram)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("Error during body close:", err)
			}
		}(r.Body)

		if err = validators.ValidateProgram(&updatedProgram); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repository.UpdateProgramByID(id, updatedProgram, env.Db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		updatedProgram.Id = new(uint)
		*updatedProgram.Id = id
		response := map[string]interface{}{
			"program": updatedProgram,
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

func (env *ProgramHandler) DeleteProgramHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing program ID", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.Atoi(idStr)
		id := uint(idInt)
		if err != nil {
			http.Error(w, "Invalid program ID", http.StatusBadRequest)
			return
		}
		_, err = repository.GetProgramByID(id, env.Db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Program not found", http.StatusNotFound)
			} else {
				log.Printf("Error fetching program: %v", err)
				http.Error(w, "Failed to fetch program", http.StatusInternalServerError)
			}
			return
		}
		err = repository.DeleteProgram(id, env.Db)
		if err != nil {
			log.Printf("Error deleting program: %v", err)
			http.Error(w, "Failed to delete program", http.StatusInternalServerError)
			return
		}
		msg := fmt.Sprintf("Program with id %d deleted successfully", id)
		response := map[string]interface{}{
			"message": msg,
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during encoding:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
