package routes

import (
	"net/http"
	"openprogramschedule/api/handlers"
)

func ProgramRouter(router *http.ServeMux, env *handlers.ProgramHandler) {
	router.HandleFunc("POST /programs/add", env.AddProgramHandler)
	router.HandleFunc("GET /programs/get-by-id", env.GetProgramByIDHandler)              // /programs/get-by-id?id
	router.HandleFunc("GET /programs/get-by-name", env.GetProgramByNameHandler)          // /programs/get-by-name?name
	router.HandleFunc("GET /programs/get-by-category", env.GetProgramsByCategoryHandler) // /programs/get-by-category?category
	router.HandleFunc("GET /programs/all", env.GetAllProgramsHandler)
	router.HandleFunc("PUT /programs/update", env.UpdateProgramHandler)          // /programs/update?id
	router.HandleFunc("DELETE /programs/delete-by-id", env.DeleteProgramHandler) // /programs/delete-by-id?id
}

func ScheduleRouter(router *http.ServeMux, env *handlers.ScheduleHandler) {
	router.HandleFunc("POST /schedules/add", env.AddScheduleHandler)
	router.HandleFunc("GET /schedules/all", env.GetAllSchedulesHandler)
	router.HandleFunc("GET /schedules/get-by-id", env.GetScheduleByIDHandler)                // /schedules/get-by-id?id
	router.HandleFunc("GET /schedules/get-by-program-id", env.GetScheduleByProgramIdHandler) // /schedules/get-by-program-id?programId
	router.HandleFunc("GET /schedules/get-by-day", env.GetScheduleByDayHandler)              // /schedules/get-by-day?day
	router.HandleFunc("GET /schedules/get-by-date", env.GetScheduleByDateHandler)            // /schedules/get-by-date?date
	router.HandleFunc("PUT /schedules/update", env.UpdateScheduleHandler)                    // /schedules/update?id
	router.HandleFunc("DELETE /schedules/delete-by-id", env.DeleteScheduleHandler)           // /schedules/delete-by-id?id
	router.HandleFunc("DELETE /schedules/delete-all", env.DeleteAllSchedulesHandler)
}
